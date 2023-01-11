package rpc

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/cors"
)

const (
	maxRequestContentLength = 1024 * 1024 * 5
	contentType             = "application/json"
)

var acceptedContentTypes = []string{contentType, "application/json-rpc", "application/jsonrequest"}

type Server struct {
	allowedOrigins []string
	vhosts         map[string]struct{}
	jwtSecret      []byte

	services *serviceRegistry
	idgen    func() ID
}

func NewServer(cors []string, vhosts []string, jwtSecret []byte, services *serviceRegistry) *Server {
	vhostMap := make(map[string]struct{})
	for _, allowedHost := range vhosts {
		vhostMap[strings.ToLower(allowedHost)] = struct{}{}
	}

	return &Server{
		allowedOrigins: cors,
		vhosts:         vhostMap,
		jwtSecret:      jwtSecret,
		services:       services,
	}
}

// ServeHTTP serves JSON-RPC requests over HTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out := w

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()
		out = &gzipResponseWriter{ResponseWriter: w, Writer: gz}
	}

	if err := s.authentication(r); err != nil {
		http.Error(out, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := s.verifyVirtualHost(r); err != nil {
		http.Error(out, err.Error(), http.StatusUnauthorized)
		return
	}

	if len(s.allowedOrigins) != 0 {
		c := cors.New(cors.Options{
			AllowedOrigins: s.allowedOrigins,
			AllowedMethods: []string{http.MethodPost, http.MethodGet},
			AllowedHeaders: []string{"*"},
			MaxAge:         600,
		})

		c.Handler(s.processRequestHandle(out))
		return
	}

	s.processRequest(out, r)
}

func (s *Server) authentication(r *http.Request) error {
	if len(s.jwtSecret) == 0 {
		return nil
	}

	var (
		strToken string
		claims   jwt.RegisteredClaims
	)
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		strToken = strings.TrimPrefix(auth, "Bearer ")
	}
	if len(strToken) == 0 {
		return errors.New("missing token")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	}

	token, err := jwt.ParseWithClaims(strToken, &claims, keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithoutClaimsValidation())

	switch {
	case err != nil:
		return err
	case !token.Valid:
		return errors.New("invalid token")
	case !claims.VerifyExpiresAt(time.Now(), false): // optional
		return errors.New("token is expired")
	case claims.IssuedAt == nil:
		return errors.New("missing issued-at")
	case time.Since(claims.IssuedAt.Time) > jwtExpiryTimeout:
		return errors.New("stale token")
	case time.Until(claims.IssuedAt.Time) > jwtExpiryTimeout:
		return errors.New("future token")
	}

	return nil
}

func (s *Server) verifyVirtualHost(r *http.Request) error {
	if r.Host == "" {
		return nil
	}

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		// Either invalid (too many colons) or no port specified
		host = r.Host
	}

	if ipAddr := net.ParseIP(host); ipAddr != nil {
		return nil
	}

	// Not an IP address, but a hostname. Need to validate
	if _, exist := s.vhosts["*"]; exist {
		return nil
	}
	if _, exist := s.vhosts[host]; exist {
		return nil
	}

	return errors.New("invalid host specified")
}

func (s *Server) processRequestHandle(out http.ResponseWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.processRequest(out, r)
	})
}

func (s *Server) processRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.ContentLength == 0 && r.URL.RawQuery == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if code, err := validateRequest(r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Create request-scoped context.
	connInfo := PeerInfo{Transport: "http", RemoteAddr: r.RemoteAddr}
	connInfo.HTTP.Version = r.Proto
	connInfo.HTTP.Host = r.Host
	connInfo.HTTP.Origin = r.Header.Get("Origin")
	connInfo.HTTP.UserAgent = r.Header.Get("User-Agent")
	ctx := r.Context()
	ctx = context.WithValue(ctx, peerInfoContextKey{}, connInfo)

	// All checks passed, create a codec that reads directly from the request body
	// until EOF, writes the response to w, and orders the server to process a
	// single request.
	w.Header().Set("content-type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	conn := NewHttpConn(r, w)
	defer conn.close()
	s.serveSingleRequest(ctx, conn)

}

func (s *Server) serveSingleRequest(ctx context.Context, conn *httpConn) {
	h := newHandler(ctx, conn, s.idgen, s.services)
	reqs, batch, err := conn.readBatch()
	if err != nil {
		if err != io.EOF {
			conn.writeJSON(ctx, errorMessage(&invalidMessageError{"parse error"}))
		}
		return
	}
	if batch {
		h.handleBatch(reqs)
	} else {
		h.handleMsg(reqs[0])
	}
}

func validateRequest(r *http.Request) (int, error) {
	if r.Method == http.MethodPut || r.Method == http.MethodDelete {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}
	if r.ContentLength > maxRequestContentLength {
		err := fmt.Errorf("content length too large (%d>%d)", r.ContentLength, maxRequestContentLength)
		return http.StatusRequestEntityTooLarge, err
	}
	// Allow OPTIONS (regardless of content-type)
	if r.Method == http.MethodOptions {
		return 0, nil
	}
	// Check content-type
	if mt, _, err := mime.ParseMediaType(r.Header.Get("content-type")); err == nil {
		for _, accepted := range acceptedContentTypes {
			if accepted == mt {
				return 0, nil
			}
		}
	}
	// Invalid content-type
	err := fmt.Errorf("invalid content type, only %s is supported", contentType)
	return http.StatusUnsupportedMediaType, err
}

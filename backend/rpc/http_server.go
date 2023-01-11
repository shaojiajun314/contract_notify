package rpc

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

)

type HttpServer struct {
	endpoint string

	mux http.ServeMux // registered handlers go here

	mu       sync.Mutex
	server   *http.Server
	listener net.Listener // non-nil when server is running

	httpHandler atomic.Value // *rpcHandler

	ws        bool
	wsHandler atomic.Value // *rpcHandler

	services serviceRegistry
}

func NewHttpServer(endpoint string, enableWS bool) (*HttpServer, error) {
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &HttpServer{
		endpoint: endpoint,
		listener: listener,
		ws:       enableWS,
	}, nil
}

func (h *HttpServer) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.endpoint == "" || h.listener == nil {
		return errors.New("failed to new http server")
	}

	h.server = &http.Server{Handler: h}

	// go h.server.Serve(h.listener)
	h.server.Serve(h.listener)

	return nil
}

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.ws {
		// h.logger.Info("start ws")
	}

	conn := NewServer(nil, nil, nil, &h.services)
	conn.ServeHTTP(w, r)
}

func (h *HttpServer) Stop() {
	// if atomic.CompareAndSwapInt32(&s.stop, 0, 1) {
	// 	fmt.Println("RPC server shutting down")
	// }
}

func (h *HttpServer) RegisterApis(apis []API) error {
	for _, api := range apis {
		if err := h.registerName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	return nil
}

func (h *HttpServer) registerName(name string, receiver interface{}) error {
	return h.services.registerName(name, receiver)
}

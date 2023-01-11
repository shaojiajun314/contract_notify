package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type HTTPAuth func(h http.Header) error

func NewJWTAuth(jwtsecret [32]byte) HTTPAuth {
	return func(h http.Header) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": &jwt.NumericDate{Time: time.Now()},
		})
		s, err := token.SignedString(jwtsecret[:])
		if err != nil {
			return fmt.Errorf("failed to create JWT token: %w", err)
		}
		h.Set("Authorization", "Bearer "+s)
		return nil
	}
}

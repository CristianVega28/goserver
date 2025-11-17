package middleware

import (
	"fmt"
	"net/http"
)

type (
	AuthMiddleware struct{}
)

func (auth *AuthMiddleware) Jwt() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("JWT Middleware")
			f(w, r)
		}

	}
}

func (auth *AuthMiddleware) BasicAuth() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Basic Auth Middleware")
			f(w, r)
		}

	}
}

func (auth *AuthMiddleware) BearerToken() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Bearer Token Middleware")
			f(w, r)
		}

	}
}

package middleware

import "net/http"

type (
	AuthMiddleware struct{}
)

func (auth *AuthMiddleware) Jwt(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (auth *AuthMiddleware) BasicAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (auth *AuthMiddleware) BearerToken(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

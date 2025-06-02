package middleware

import "net/http"

type (
	SecurityMiddleware struct{}
)

func (security *SecurityMiddleware) Csrf(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (security *SecurityMiddleware) Cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (security *SecurityMiddleware) RateLimit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type (
	SecurityMiddleware struct{}
)

func (security *SecurityMiddleware) Csrf() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("CSRF Middleware")
			f(w, r)
		}
	}
}

func (security *SecurityMiddleware) Cors() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
}

func (security *SecurityMiddleware) RateLimit() MiddlewareFunction {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			f(w, r)
		}

	}
}

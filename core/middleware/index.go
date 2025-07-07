package middleware

import (
	"net/http"

	"github.com/CristianVega28/goserver/helpers"
)

type (
	MiddlewareFunction func(w http.HandlerFunc) http.HandlerFunc

	// helper helpers.Response{}
)

var helper helpers.Response = helpers.Response{}

func Get(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// helper.ResponseJson(w, map[string]string{}, http.StatusBadRequest)
		}

		next(w, r)
	}
}
func Post(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			helper.ResponseJson(w, map[string]string{
				"error": "No esta permitido ese tipo de solicit",
			}, http.StatusBadRequest)
		}

		next(w, r)
	}
}
func Delete(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			helper.ResponseJson(w, map[string]string{
				"error": "No esta permitido ese tipo de solicit",
			}, http.StatusBadRequest)
		}

		next(w, r)
	}
}

func Put(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			helper.ResponseJson(w, map[string]string{
				"error": "No esta permitido ese tipo de solicit",
			}, http.StatusBadRequest)
		}

		next(w, r)
	}
}

func Chain(f http.HandlerFunc, middleware ...MiddlewareFunction) http.HandlerFunc {
	for _, v := range middleware {
		v(f)
	}

	return f
}

func FunctionsAuthMiddleware() []MiddlewareFunction {
	authmid := AuthMiddleware{}
	return []MiddlewareFunction{
		authmid.BasicAuth,
		authmid.BearerToken,
		authmid.Jwt,
	}
}

func FunctionsSecurityMiddleware() []MiddlewareFunction {
	securitymid := SecurityMiddleware{}
	return []MiddlewareFunction{
		securitymid.Cors,
		securitymid.Csrf,
		securitymid.RateLimit,
	}
}

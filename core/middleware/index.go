package middleware

import "net/http"

type (
	MiddlewareFunction func(w http.HandlerFunc) http.HandlerFunc
)

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

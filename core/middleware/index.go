package middleware

import (
	"net/http"

	"github.com/CristianVega28/goserver/helpers"
)

type (
	MiddlewareFunction func(w http.HandlerFunc) http.HandlerFunc
	MapMiddleware      map[string][]MiddlewareFunction
)

func Chain(f http.HandlerFunc, middleware ...MiddlewareFunction) http.HandlerFunc {
	for _, v := range middleware {
		f = v(f)
	}

	return f
}

func FunctionsAuthMiddleware(key string) MiddlewareFunction {
	authmid := AuthMiddleware{}
	var mf MiddlewareFunction
	if key == "basic_auth" {
		mf = authmid.BasicAuth()
	}

	if key == "bearer" {
		mf = authmid.BearerToken()
	}

	if key == "jwt" {
		mf = authmid.Jwt()
	}

	return mf
}

func FunctionsSecurityMiddleware(keys []string) []MiddlewareFunction {
	securitymid := SecurityMiddleware{}
	var arraySecurity []MiddlewareFunction
	for _, v := range keys {
		if v == "cors" {
			arraySecurity = append(arraySecurity, securitymid.Cors())
		}
		if v == "csrf" {
			arraySecurity = append(arraySecurity, securitymid.Csrf())
		}
		if v == "rate_limit" {
			arraySecurity = append(arraySecurity, securitymid.RateLimit())
		}

	}

	return arraySecurity
}

func ReturnArraysMiddleware(cfg helpers.ConfigServerApi) []MiddlewareFunction {
	var arrMiddleware []MiddlewareFunction

	if cfg.MiddlewareApi.Security != nil {
		securityMiddleware := FunctionsSecurityMiddleware(cfg.MiddlewareApi.Security)
		arrMiddleware = append(arrMiddleware, securityMiddleware...)
	}

	if cfg.MiddlewareApi.Auth != "" {
		authMiddleware := FunctionsAuthMiddleware(cfg.MiddlewareApi.Auth)
		arrMiddleware = append(arrMiddleware, authMiddleware)
	}

	return arrMiddleware
}

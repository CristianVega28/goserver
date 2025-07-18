package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type (
	MiddlewareFunction func(w http.HandlerFunc) http.HandlerFunc
	MapMiddleware      map[string][]MiddlewareFunction

	// helper helpers.Response{}
)

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

func CreateMapMiddleware() MapMiddleware {
	entries := nameMiddlewareFromDirectories()
	mappingMiddleware := make(MapMiddleware)

	for _, file := range entries {
		name := strings.Split(file.Name(), ".")[0]
		switch name {
		case "auth":
			mappingMiddleware[name] = FunctionsAuthMiddleware()

		case "security":
			mappingMiddleware[name] = FunctionsSecurityMiddleware()

		}
	}
	return mappingMiddleware
}

func nameMiddlewareFromDirectories() []os.DirEntry {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	pathMiddleware := filepath.Join(cwd, "/../../core/middleware")
	filesMiddleware, _ := os.ReadDir(pathMiddleware)

	return filesMiddleware
}

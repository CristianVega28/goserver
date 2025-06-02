package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/CristianVega28/goserver/core/middleware"
)

type (
	Methods struct {
		MiddKeywords []string
		PathKeywords string
		Server       *http.ServeMux
	}

	MapMiddleware map[string][]middleware.MiddlewareFunction
)

func (method *Methods) GenerateControllers() {
	//let's create a map where it store the key (middleware) and value of which one it save the middleware function.

}

func CreateMapMiddleware() MapMiddleware {
	entries := nameMiddlewareFromDirectories()
	fmt.Println(entries)
	mappingMiddleware := make(MapMiddleware)

	for _, file := range entries {
		name := strings.Split(file.Name(), ".")[0]
		switch name {
		case "auth":
			mappingMiddleware[name] = middleware.FunctionsAuthMiddleware()

		case "security":
			mappingMiddleware[name] = middleware.FunctionsSecurityMiddleware()

		}
	}
	return mappingMiddleware
}

func nameMiddlewareFromDirectories() []os.DirEntry {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Println(cwd)
	pathMiddleware := filepath.Join(cwd, "/../../core/middleware")
	fmt.Println(pathMiddleware)

	filesMiddleware, _ := os.ReadDir(pathMiddleware)

	return filesMiddleware
}

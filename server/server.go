package server

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/CristianVega28/goserver/core/middleware"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/rs/zerolog/log"
)

type (
	Server struct {
		mux *http.ServeMux
		Srv http.Server
	}
)

func (server *Server) NewServer() Server {
	return Server{
		mux: http.NewServeMux(),
		Srv: http.Server{},
	}
}

func (server *Server) Up(serverVar *http.Server) {

	err := serverVar.ListenAndServe()

	if err != nil {
		log.Error().Msg("Not working server")
	}

}

func (server *Server) Close() {

	fmt.Println("\nApagando servidor...")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	if err := server.Srv.Close(); err != nil {
		fmt.Printf("Error cerrando servidor: %s\n", err)
	} else {
		fmt.Println("Servidor cerrado correctamente.")
	}

	fmt.Println("Reiniciando servidor...")
}

func (server *Server) GenrateServer(data map[string]any) {

	response := helpers.Response{}

	if len(data) != 0 {
		for key, value := range data {

			typeValue := reflect.TypeOf(value)
			fmt.Println(typeValue.Kind())
			switch typeValue.Kind() {
			case reflect.Slice:

				funcWithMiddleware := middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					response := helpers.Response{}
					response.ResponseJson(w, value, http.StatusOK)

				}, middleware.Logging)

				// it create GET, POST , DELETE, PUT
				path := fmt.Sprintf("/%s", key)
				server.mux.HandleFunc(path, middleware.Get(funcWithMiddleware))
				// server.mux.HandleFunc(path, middleware.Post(funcWithMiddleware))
				// server.mux.HandleFunc(path, middleware.Delete(funcWithMiddleware))
				// server.mux.HandleFunc(path, middleware.Put(funcWithMiddleware))

			case reflect.Map:
				fmt.Println(value)

			}

			// server.mux.HandleFunc("/"+key, middleware.Chain(func(w http.ResponseWriter, r *http.Request){} , middleware.Logging))
		}
	}
	server.mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		response.ResponseJson(w, map[string]any{
			"code":    http.StatusAccepted,
			"message": "up",
		}, http.StatusAccepted)
	})

	server.Srv.Handler = server.mux

	server.Srv.ListenAndServe()
}
func GeneralFunc(value any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

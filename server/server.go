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

func (server *Server) GenrateServer(data map[string]any) {

	response := helpers.Response{}

	if len(data) != 0 {
		for key, value := range data {

			typeValue := reflect.TypeOf(value)
			switch typeValue.Kind() {
			case reflect.Slice:
				funcWithMiddleware := middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {

					case http.MethodGet:
						Get(w, r, nil)
					case http.MethodPost:
						Post(w, r, nil)
					case http.MethodDelete:
						Delete(w, r, nil)
					case http.MethodPut:
						Put(w, r, nil)

					}
				}, middleware.Logging)

				// it create GET, POST , DELETE, PUT
				path := fmt.Sprintf("/%s", key)
				server.mux.HandleFunc(path, funcWithMiddleware)
				// server.mux.HandleFunc(path, middleware.Post(funcWithMiddleware))
				// server.mux.HandleFunc(path, middleware.Delete(funcWithMiddleware))
				// server.mux.HandleFunc(path, middleware.Put(funcWithMiddleware))

			case reflect.Map:
			}

			// server.mux.HandleFunc("/"+key, middleware.Chain(func(w http.ResponseWriter, r *http.Request){} , middleware.Logging))
		}
	}
	server.mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		response.ResponseJson(w, map[string]any{
			"code":    http.StatusAccepted,
			"message": "up change",
		}, http.StatusAccepted)
	})

	server.Srv.Handler = server.mux

	server.Srv.ListenAndServe()
}
func GeneralFunc(value any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

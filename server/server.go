package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/CristianVega28/goserver/core/db"
	"github.com/CristianVega28/goserver/core/middleware"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/CristianVega28/goserver/utils"
)

type (
	Server struct {
		mux *http.ServeMux
		Srv http.Server
	}
)

var logs = utils.Logger{}
var log = logs.Create()

func (server *Server) NewServer() Server {
	return Server{
		mux: http.NewServeMux(),
		Srv: http.Server{},
	}
}

func (server *Server) GenrateServer(data map[string]any) {

	response := helpers.Response{}

	if len(data) != 0 {
		for key, value := range data {

			typeValue := reflect.TypeOf(value)
			path := fmt.Sprintf("/%s", key)
			switch typeValue.Kind() {
			case reflect.Slice:
				// it create GET, POST , DELETE, PUT
				server.mux.HandleFunc(path, middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						Get(w, r, nil, value)
					case http.MethodPost:
						Post(w, r, nil)
					case http.MethodDelete:
						Delete(w, r, nil)
					case http.MethodPut:
						Put(w, r, nil)
					}
				}, middleware.Logging()))
			case reflect.Map:
				var cfg helpers.ConfigServerApi
				valueResponse, err := json.Marshal(value)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(valueResponse, &cfg)
				if err != nil {
					fmt.Println("Error unmarshalling config:", err)
					continue
				}

				var funcRequest http.HandlerFunc

				funcRequest = middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						Get(w, r, &cfg, cfg.Response)
					case http.MethodPost:
						Post(w, r, &cfg)
					case http.MethodDelete:
						Delete(w, r, &cfg)
					case http.MethodPut:
						Put(w, r, &cfg)
					}

				})

				// log.Structs("Configuracion cfg", cfg)
				SetConfigurationServer(cfg)
				server.mux.HandleFunc(path, funcRequest)

			}
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

func SetConfigurationServer(cfg helpers.ConfigServerApi) {

	if cfg.Schema != nil {
		log.Structs("Schema", cfg.Schema)
		db.MigrateSchema(cfg.Schema)

	}
}

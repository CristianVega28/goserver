package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/CristianVega28/goserver/core/middleware"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/rs/cors"
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

func (server *Server) GenrateServer(data map[string]any) {

	response := helpers.Response{}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // o ej: []string{"http://localhost:4321"}
		// AllowedMethods: []string{"*"},
		// AllowedHeaders:   []string{"Content-Type", "Authorization"},
		// AllowCredentials: true,
	})
	var arrCfgResponse []helpers.ResponseConfig

	if len(data) != 0 {
		for key, value := range data {

			typeValue := reflect.TypeOf(value)
			path := fmt.Sprintf("/%s", key)
			switch typeValue.Kind() {
			case reflect.Slice:
				// it create GET, POST , DELETE, PUT
				rspCfg := helpers.ResponseConfig{
					Path:    path,
					Request: []string{"GET", "POST", "DELETE", "PUT"},
				}
				arrCfgResponse = append(arrCfgResponse, rspCfg)
				server.mux.HandleFunc(path, middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						Get(w, r, value)
					case http.MethodPost:
						Post(w, r)
					case http.MethodDelete:
						Delete(w, r)
					case http.MethodPut:
						Put(w, r)
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

				rspCfg := helpers.ResponseConfig{
					Path:    path,
					Request: cfg.Request,
					Schema:  cfg.Schema,
				}

				arrCfgResponse = append(arrCfgResponse, rspCfg)
				var funcRequest http.HandlerFunc

				funcRequest = middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), helpers.KeyCfg, cfg)
					switch r.Method {
					case http.MethodGet:
						Get(w, r.WithContext(ctx), cfg.Response)
					case http.MethodPost:
						Post(w, r.WithContext(ctx))
					case http.MethodDelete:
						Delete(w, r.WithContext(ctx))
					case http.MethodPut:
						Put(w, r.WithContext(ctx))
					}

				})

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

	server.mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	server.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public/"))))

	server.mux.HandleFunc("/docs-api", func(w http.ResponseWriter, r *http.Request) {

		response.ResponseJson(w, map[string]any{
			"success": true,
			"data":    arrCfgResponse,
		}, http.StatusOK)
	})

	handler := c.Handler(server.mux)
	server.Srv.Handler = handler

	server.Srv.ListenAndServe()
}

func SetConfigurationServer(cfg helpers.ConfigServerApi) {

	if cfg.Schema != nil {

		// Here create the tables in database
		model := helpers.MigrateSchema(cfg.Schema)
		response := checkTypesForResponse(cfg.Response)
		model.InsertMigration(response, false)

	}
}

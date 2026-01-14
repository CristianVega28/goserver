package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/CristianVega28/goserver/core/controllers"
	"github.com/CristianVega28/goserver/core/middleware"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/rs/cors"
)

type (
	Server struct {
		mux   *http.ServeMux
		Srv   http.Server
		Debug bool
	}
)

func (server *Server) NewServer() Server {
	return Server{
		mux:   http.NewServeMux(),
		Srv:   http.Server{},
		Debug: server.Debug,
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

	var statistics helpers.ConfigServerStatistics = helpers.ConfigServerStatistics{}

	(&statistics).Loader(data)

	statistics.TotalRecords = 0

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

				if _v, ok := value.([]any); ok {
					statistics.TotalRecords += len(_v)
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

				arrMiddleware := middleware.ReturnArraysMiddleware(cfg)

				fmt.Println(cfg.MiddlewareApi.Auth)
				if _value, ok := cfg.Response.([]any); ok {
					statistics.TotalRecords += len(_value)
				}

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

				}, arrMiddleware...)

				SetConfigurationServer(cfg)

				if cfg.MiddlewareApi.Auth == "bearer" {
					pathTemp := fmt.Sprintf("/%s/token", path)
					server.mux.HandleFunc(pathTemp, (&controllers.AuthController{}).BearerController())
				}

				server.mux.HandleFunc(path, funcRequest)

			}
		}
	}
	server.mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		response.ResponseJson(w, map[string]any{
			"code":    http.StatusAccepted,
			"message": "up change jsjs",
		}, http.StatusAccepted)
	})

	server.mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	server.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public/"))))
	// server.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// http.ServeFile(w, r, "public/index.html")
	// 	http.Redirect(w, r, "/docs", http.StatusSeeOther)
	// })

	server.mux.HandleFunc("/docs-api", func(w http.ResponseWriter, r *http.Request) {
		response.ResponseJson(w, map[string]any{
			"success": true,
			"data":    arrCfgResponse,
		}, http.StatusOK)
	})

	server.mux.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		response.ResponseJson(w, map[string]any{
			"success": true,
			"data": map[string]any{
				"total_requests": statistics.TotalRequests,
				"total_tables":   statistics.TotalTables,
				"total_records":  statistics.TotalRecords,
			},
		}, http.StatusOK)
	})

	handler := c.Handler(server.mux)
	server.Srv.Handler = handler

	if !server.Debug {
		server.Srv.ListenAndServe()
	}
}

func SetConfigurationServer(cfg helpers.ConfigServerApi) {

	cfg.PreLoader()

	if cfg.Schema != nil {

		// Here create the tables in database
		model := helpers.MigrateSchema(cfg.Schema)
		model.SetResponse(cfg.Response)
		model.InsertMigration(false)

	}
}

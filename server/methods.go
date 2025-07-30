package server

import (
	"net/http"

	"github.com/CristianVega28/goserver/helpers"
	"github.com/samber/lo"
)

type (
	Methods struct {
		MiddKeywords []string
		PathKeywords string
		Server       *http.ServeMux
	}
)

var helper helpers.Response = helpers.Response{}

func Get(w http.ResponseWriter, r *http.Request, cfg *helpers.ConfigServerApi, values any) error {

	if cfg == nil {
		helper.ResponseJson(w, values, http.StatusAccepted)
		return nil
	}
	if valid := ValidationCfgMethod(r.Method, cfg.Request); !valid {
		helper.ResponseJson(w, map[string]string{
			"validated": "Method not allowed",
		}, http.StatusMethodNotAllowed)
		return nil
	}

	helper.ResponseJson(w, values, http.StatusAccepted)
	return nil
}
func Post(w http.ResponseWriter, r *http.Request, cfg *helpers.ConfigServerApi) error {
	if cfg == nil {
		return nil
	}
	if valid := ValidationCfgMethod(r.Method, cfg.Request); !valid {
		helper.ResponseJson(w, map[string]string{
			"validated": "Method not allowed",
		}, http.StatusMethodNotAllowed)
		return nil
	}
	return nil
}
func Delete(w http.ResponseWriter, r *http.Request, cfg *helpers.ConfigServerApi) error {
	if cfg == nil {
		return nil
	}
	if valid := ValidationCfgMethod(r.Method, cfg.Request); !valid {
		helper.ResponseJson(w, map[string]string{
			"validated": "Method not allowed",
		}, http.StatusMethodNotAllowed)
		return nil
	}
	return nil
}
func Put(w http.ResponseWriter, r *http.Request, cfg *helpers.ConfigServerApi) error {
	if cfg == nil {
		return nil
	}
	if valid := ValidationCfgMethod(r.Method, cfg.Request); !valid {
		helper.ResponseJson(w, map[string]string{
			"validated": "Method not allowed",
		}, http.StatusMethodNotAllowed)
		return nil
	}
	return nil
}

func ValidationCfgMethod(method string, methods []string) bool {
	include := lo.Contains(methods, method)
	return include
}

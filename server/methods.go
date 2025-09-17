package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/samber/lo"
)

type (
	Methods struct {
		MiddKeywords []string
		PathKeywords string
		Server       *http.ServeMux
	}
	StructPost struct {
		Models any `json:"models"`
	}
)

var helper helpers.Response = helpers.Response{}

func Get(w http.ResponseWriter, r *http.Request, values any) error {

	cfg, ok := r.Context().Value(helpers.KeyCfg).(helpers.ConfigServerApi)
	if !ok {
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
func Post(w http.ResponseWriter, r *http.Request) error {
	cfg, ok := r.Context().Value(helpers.KeyCfg).(helpers.ConfigServerApi)

	if !ok {
		helper.ResponseJson(w, map[string]any{
			"success": false,
		}, http.StatusInternalServerError)
		return nil
	}
	if valid := ValidationCfgMethod(r.Method, cfg.Request); !valid {
		helper.ResponseJson(w, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		}, http.StatusMethodNotAllowed)
		return nil
	}

	modelBk := models.Models[map[string]any]{}
	model := modelBk.Init()
	metadata := cfg.ReturnMetadataTable()
	model.SetMetadataTable(metadata)
	model.SetTableName(cfg.Schema["table_name"].(string))
	var body StructPost
	err := json.NewDecoder(r.Body).Decode(&body)

	if err == io.EOF {
		helper.ResponseJson(w, map[string]any{
			"success": false,
			"error":   "body is empty",
		}, http.StatusBadRequest)
		return nil
	}

	errors := model.ValidateFields(body.Models)

	if len(errors) > 0 {
		helper.ResponseJson(w, errors, http.StatusUnprocessableEntity)
		return nil
	}
	response := checkTypesForResponse(body.Models)
	errInsert := model.InsertMigration(response, true)

	if errInsert != nil {
		helper.ResponseJson(w, map[string]string{
			"error": errInsert.Error(),
		}, http.StatusBadRequest)
		return nil
	}
	description := fmt.Sprintf("Inserted %d rows", len(response))

	helper.ResponseJson(w, map[string]any{
		"success": true,
		"status":  description,
	}, http.StatusCreated)
	return nil
}
func Delete(w http.ResponseWriter, r *http.Request) error {
	cfg, ok := r.Context().Value(helpers.KeyCfg).(helpers.ConfigServerApi)
	if !ok {
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
func Put(w http.ResponseWriter, r *http.Request) error {
	cfg, ok := r.Context().Value(helpers.KeyCfg).(helpers.ConfigServerApi)
	if !ok {
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

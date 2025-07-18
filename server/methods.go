package server

import (
	"net/http"

	"github.com/CristianVega28/goserver/helpers"
)

type (
	Methods struct {
		MiddKeywords []string
		PathKeywords string
		Server       *http.ServeMux
	}
)

var helper helpers.Response = helpers.Response{}

func Get(w http.ResponseWriter, r *http.Request, cfg any) error {
	return nil
}
func Post(w http.ResponseWriter, r *http.Request, cfg any) error {
	return nil
}
func Delete(w http.ResponseWriter, r *http.Request, cfg any) error {
	return nil
}
func Put(w http.ResponseWriter, r *http.Request, cfg any) error {
	return nil
}

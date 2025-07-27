package middleware

import (
	"fmt"
	"net/http"

	"github.com/CristianVega28/goserver/utils"
)

func Logging() func(http.HandlerFunc) http.HandlerFunc {

	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logs := utils.Logger{}
			log := logs.Create()
			log.Msg(fmt.Sprintf("Method: %s, Path: %s", r.Method, r.URL.Path))
			f(w, r)
		}

	}
}

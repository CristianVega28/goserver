package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/CristianVega28/goserver/utils"
)

type (
	AuthMiddleware struct{}
)

const (
	nameCookie = "sessionid"
)

var helper helpers.Response = helpers.Response{}

func (auth *AuthMiddleware) Jwt() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("JWT Middleware")
			f(w, r)
		}

	}
}

func (auth *AuthMiddleware) BasicAuth() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Basic Auth Middleware")
			f(w, r)
		}

	}
}

func (auth *AuthMiddleware) BearerToken() func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(nameCookie)
			if err != nil {
				if err == http.ErrNoCookie {
					helper.ResponseJson(w, map[string]any{
						"success": false,
						"errors":  "Not authenticated",
					}, http.StatusUnauthorized)

				}
				return
			}
			token, exist := models.Cache_.Get(cookie.Value)

			if !exist {
				helper.ResponseJson(w, map[string]any{
					"success": false,
					"errors":  "Not authenticated",
				}, http.StatusUnauthorized)
				return
			}

			auth := r.Header.Get("Authorization")
			_tokenR, existBearer := utils.GetAt(strings.Fields(auth), 1)
			_tokenM := fmt.Sprint(token.Value)

			if !existBearer || _tokenR != _tokenM {
				helper.ResponseJson(w, map[string]any{
					"success": false,
					"errors":  "Not authenticated",
				}, http.StatusUnauthorized)
				return

			}

			f(w, r)
		}

	}
}

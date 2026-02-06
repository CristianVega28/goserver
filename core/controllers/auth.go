package controllers

import (
	"net/http"
	"strconv"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/helpers"
)

type (
	AuthController struct{}
)

const (
	nameCookie = "sessionid"
)

var helper helpers.Response = helpers.Response{}

func (auth *AuthController) BearerController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(nameCookie)
		if err != nil {
			if err == http.ErrNoCookie {

				token, expiration := helper.GenerateToken(w)

				helper.ResponseJson(w, map[string]string{
					"api_key": token,
					"expires": expiration,
				}, http.StatusOK)

			}
			return
		}

		token, exist := models.Cache_.Get(cookie.Value)

		if exist == false {

			helper.DeleteSetCookie(w, nameCookie)
			token, expiration := helper.GenerateToken(w)

			helper.ResponseJson(w, map[string]string{
				"error":       "Not found token",
				"detail":      "The token associated with the session cookie was not found or has expired",
				"new_api_key": token,
				"expires":     expiration,
			}, http.StatusNotFound)
			return
		}

		helper.ResponseJson(w, map[string]string{
			"api_key": token.Value.(string),
			"expires": strconv.FormatInt(token.Ttl, 10),
		}, http.StatusOK)

		return

	}

}

func (auth *AuthController) GetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionid, err := r.Cookie("sessionid")

		if err != nil {
			if err == http.ErrNoCookie {
				helper.ResponseJson(w, map[string]string{
					"error": "No session cookie found",
				}, http.StatusUnauthorized)
				return
			}
		}

		value, exist := models.Cache_.Get(sessionid.Value)

		if exist == false {
			helper.ResponseJson(w, map[string]string{
				"error": "Not found token",
			}, http.StatusNotFound)

		}

		helper.ResponseJson(w, map[string]string{
			"token": value.Value.(string),
		}, http.StatusOK)
	}
}

package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/helpers"
	"github.com/CristianVega28/goserver/utils"
	"github.com/google/uuid"
)

type (
	AuthController struct{}
)

var helper helpers.Response = helpers.Response{}

func (auth *AuthController) BearerController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes := make([]byte, 64)
		if _, err := rand.Read(bytes); err != nil {
			utils.Log.Fatal(err.Error())
		}
		cookie, err := r.Cookie("sessionid")

		if err != nil {
			if err == http.ErrNoCookie {

				id := uuid.New().String()

				http.SetCookie(w, &http.Cookie{
					Name:     "sessionid",
					Value:    id,
					HttpOnly: true,
				})
				token := hex.EncodeToString(bytes)
				expiration := time.Now().Unix()

				models.Cache_.Set(id, token, expiration)

				helper.ResponseJson(w, map[string]string{
					"api_key": token,
					"expires": strconv.FormatInt(expiration, 10),
				}, http.StatusOK)

			}
		}

		token, bool := models.Cache_.Get(cookie.Value)

		if bool == false {
			helper.ResponseJson(w, map[string]string{
				"error": "Not found token",
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

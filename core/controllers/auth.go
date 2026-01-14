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

		var auth models.Auth = models.Auth{}

		toekn := hex.EncodeToString(bytes)
		expiration := time.Now().Unix() + int64(auth.GetBearerTokenExpiration())

		helper.ResponseJson(w, map[string]string{
			"api_key": toekn,
			"expires": strconv.FormatInt(expiration, 10),
		}, http.StatusOK)

	}

}

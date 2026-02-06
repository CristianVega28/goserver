package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/utils"
	"github.com/google/uuid"
)

type (
	Response struct{}
	Request  struct{}

	ErrorResponse struct {
		Status     int         `json:"status"`
		ErrorField interface{} `json:"error"`
	}
)

const (
	ContentTypeKey  = "Content-type"
	ApplicationJson = "application/json; charset=utf-8"
)

func (response *Response) ResponseJson(r http.ResponseWriter, information any, code int) {
	r.Header().Add(ContentTypeKey, ApplicationJson)
	r.WriteHeader(code)

	output, err := json.Marshal(information)

	if err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		badJson, _ := json.Marshal(ErrorResponse{
			Status:     http.StatusInternalServerError,
			ErrorField: err.Error(),
		})

		r.Write(badJson)
	}

	r.Write(output)

}

func (reqsponse *Response) GenerateToken(r http.ResponseWriter) (string, string) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		utils.Log.Fatal(err.Error())
	}
	id := uuid.New().String()

	http.SetCookie(r, &http.Cookie{
		Name:     "sessionid",
		Value:    id,
		HttpOnly: true,
	})
	token := hex.EncodeToString(bytes)
	expiration := time.Now().Unix()

	models.Cache_.Set(id, token, expiration)

	return token, strconv.FormatInt(expiration, 10)

}

func (response *Response) DeleteSetCookie(r http.ResponseWriter, name string) {
	http.SetCookie(r, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

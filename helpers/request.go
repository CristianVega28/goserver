package helpers

import (
	"encoding/json"
	"net/http"
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

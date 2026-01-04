package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpRouteServer(t *testing.T) {
	prevServer := Server{
		Debug: true,
	}
	srv := (&prevServer).NewServer()
	srv.GenrateServer(nil)

	req := httptest.NewRequest(http.MethodGet, "/up", nil)
	res := httptest.NewRecorder()
	srv.Srv.Handler.ServeHTTP(res, req)

	if res.Code != http.StatusAccepted {
		t.Fatal("Expected status code 200, got ", res.Code)
	}

}

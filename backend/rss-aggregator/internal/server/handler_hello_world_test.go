package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestHelloWorldHandler(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := gin.New()
	router.GET("/", apiCfg.HelloWorldHandler())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Errorf("status = %d, want %d", got, want)
	}

	if got, want := recorder.Body.String(), `{"message":"Hello World"}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

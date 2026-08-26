package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// The validation runs before the query, so these cases need no database.
func TestCreateFeedHandlerValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty body", body: `{}`},
		{name: "missing url", body: `{"name":"Nyaa"}`},
		{name: "missing name", body: `{"url":"https://nyaa.si/?page=rss"}`},
		{name: "invalid json", body: `{`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiCfg := &ApiConfig{}

			router := gin.New()
			router.POST("/feeds", apiCfg.CreateFeedHandler())

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/feeds", strings.NewReader(tt.body)))

			if got, want := recorder.Code, http.StatusBadRequest; got != want {
				t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
			}
		})
	}
}

func TestUpdateFeedHandlerValidation(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := gin.New()
	router.PUT("/feeds/:id", apiCfg.UpdateFeedHandler())

	recorder := httptest.NewRecorder()
	body := strings.NewReader(`{"url":"https://nyaa.si/?page=rss"}`)
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPut, "/feeds/6a9d1f0e-9e0e-4a4e-9a3b-2d3f4b5c6d7e", body))

	if got, want := recorder.Code, http.StatusBadRequest; got != want {
		t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
	}
}

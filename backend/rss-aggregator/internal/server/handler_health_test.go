package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthHandlerWithoutDatabase pins the behaviour that used to terminate
// the whole process: an unreachable database must be reported, not fatal.
func TestHealthHandlerWithoutDatabase(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := gin.New()
	router.GET("/health", apiCfg.HealthHandler())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	router.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusServiceUnavailable; got != want {
		t.Errorf("status = %d, want %d", got, want)
	}
}

func TestHealthWithoutDatabase(t *testing.T) {
	stats := Health(nil)

	if got, want := stats["status"], "down"; got != want {
		t.Errorf("status = %q, want %q", got, want)
	}

	if stats["error"] == "" {
		t.Error("error = \"\", want a reason for the failure")
	}
}

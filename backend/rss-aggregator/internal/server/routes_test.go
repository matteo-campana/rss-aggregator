package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	apiCfg := &ApiConfig{}

	handler, ok := apiCfg.RegisterRoutes().(*gin.Engine)
	if !ok {
		t.Fatal("RegisterRoutes did not return a *gin.Engine")
	}

	registered := map[string]bool{}
	for _, route := range handler.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	want := []string{
		http.MethodGet + " /api/v1/",
		http.MethodGet + " /api/v1/health",
		http.MethodPost + " /api/v1/users/",
		http.MethodGet + " /api/v1/users/:id",
		http.MethodGet + " /api/v1/users/me",
		http.MethodGet + " /api/v1/feeds/",
		http.MethodPost + " /api/v1/feeds/",
		http.MethodGet + " /api/v1/feed-follows/",
		http.MethodGet + " /api/v1/feed-follows/user/:user_id",
		http.MethodGet + " /api/v1/nyaa/rss",
		http.MethodGet + " /api/v1/items/",
		http.MethodGet + " /api/v1/items/categories",
		http.MethodGet + " /api/v1/channels/",
	}

	for _, route := range want {
		if !registered[route] {
			t.Errorf("route %q is not registered", route)
		}
	}
}

// The public surface is deliberately small: everything else needs an API key.
func TestPublicRoutesAreReachableWithoutAnApiKey(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := apiCfg.RegisterRoutes()

	for _, path := range []string{"/api/v1/", "/api/v1/health"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

		if recorder.Code == http.StatusUnauthorized {
			t.Errorf("%s answered 401, want it to stay public", path)
		}
	}
}

func TestProtectedRoutesRejectAnonymousRequests(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := apiCfg.RegisterRoutes()

	protected := []string{
		"/api/v1/users/",
		"/api/v1/users/me",
		"/api/v1/feeds/",
		"/api/v1/feed-follows/",
		"/api/v1/nyaa/rss",
		"/api/v1/items/",
		"/api/v1/items/categories",
		"/api/v1/channels/",
	}

	for _, path := range protected {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

		if got, want := recorder.Code, http.StatusUnauthorized; got != want {
			t.Errorf("%s status = %d, want %d", path, got, want)
		}
	}
}

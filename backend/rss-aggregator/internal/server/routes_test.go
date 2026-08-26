package server

import (
	"net/http"
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
		http.MethodGet + " /api/v1/feeds/",
		http.MethodPost + " /api/v1/feeds/",
		http.MethodGet + " /api/v1/feed-follows/",
		http.MethodGet + " /api/v1/feed-follows/user/:user_id",
		http.MethodGet + " /api/v1/nyaa/rss",
	}

	for _, route := range want {
		if !registered[route] {
			t.Errorf("route %q is not registered", route)
		}
	}
}

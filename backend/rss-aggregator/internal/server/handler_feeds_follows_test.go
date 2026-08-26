package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rss-aggregator/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// withUser mimics what MiddlewareAuth does, without a database behind it.
func withUser(user database.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(userContextKey, user)
		c.Next()
	}
}

// The handlers must take the user from the request context. Without one they
// have no business guessing, so they refuse rather than falling back to a
// value from the path or the body.
func TestFeedFollowHandlersRequireAnAuthenticatedUser(t *testing.T) {
	apiCfg := &ApiConfig{}

	tests := []struct {
		name    string
		method  string
		path    string
		body    string
		handler gin.HandlerFunc
	}{
		{"create", http.MethodPost, "/feed-follows/", `{"feed_id":"6a9d1f0e-9e0e-4a4e-9a3b-2d3f4b5c6d7e"}`, apiCfg.CreateFeedFollowsHandler()},
		{"list", http.MethodGet, "/feed-follows/", "", apiCfg.GetFeedsFollowsHandler()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Handle(tt.method, tt.path, tt.handler)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body)))

			if got, want := recorder.Code, http.StatusUnauthorized; got != want {
				t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
			}
		})
	}
}

// The body no longer carries user_id, and feed_id is still required. Both
// checks run before the query, so no database is needed.
func TestCreateFeedFollowValidation(t *testing.T) {
	apiCfg := &ApiConfig{}
	user := database.User{ID: uuid.New(), Fullname: "Matteo"}

	for _, body := range []string{`{}`, `{"user_id":"6a9d1f0e-9e0e-4a4e-9a3b-2d3f4b5c6d7e"}`, `{`} {
		router := gin.New()
		router.POST("/feed-follows/", withUser(user), apiCfg.CreateFeedFollowsHandler())

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/feed-follows/", strings.NewReader(body)))

		if got, want := recorder.Code, http.StatusBadRequest; got != want {
			t.Errorf("body %s: status = %d, want %d", body, got, want)
		}
	}
}

func TestFeedFollowHandlersRejectMalformedIds(t *testing.T) {
	apiCfg := &ApiConfig{}
	user := database.User{ID: uuid.New()}

	tests := []struct {
		name    string
		method  string
		route   string
		path    string
		handler gin.HandlerFunc
	}{
		{"get by id", http.MethodGet, "/feed-follows/:id", "/feed-follows/not-a-uuid", apiCfg.GetFeedFollowsByIdHandler()},
		{"delete by id", http.MethodDelete, "/feed-follows/:id", "/feed-follows/not-a-uuid", apiCfg.DeleteFeedFollowsHandler()},
		{"get by feed", http.MethodGet, "/feed-follows/feed/:feed_id", "/feed-follows/feed/nope", apiCfg.GetFeedFollowsByFeedIdHandler()},
		{"delete by feed", http.MethodDelete, "/feed-follows/feed/:feed_id", "/feed-follows/feed/nope", apiCfg.DeleteFeedFollowsByFeedIdHandler()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Handle(tt.method, tt.route, withUser(user), tt.handler)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(tt.method, tt.path, nil))

			if got, want := recorder.Code, http.StatusBadRequest; got != want {
				t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
			}
		})
	}
}

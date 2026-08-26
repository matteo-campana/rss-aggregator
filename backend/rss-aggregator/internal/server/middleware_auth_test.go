package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"rss-aggregator/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeUserStore struct {
	user database.User
	err  error
}

func (f fakeUserStore) GetUserByApiKey(ctx context.Context, apiKey string) (database.User, error) {
	if f.err != nil {
		return database.User{}, f.err
	}

	if apiKey != f.user.ApiKey {
		return database.User{}, sql.ErrNoRows
	}

	return f.user, nil
}

func newAuthenticatedRouter(store UserStore, reached *bool) *gin.Engine {
	router := gin.New()
	router.GET("/protected", MiddlewareAuth(store), func(c *gin.Context) {
		*reached = true

		user, ok := UserFromContext(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"fullname": user.Fullname})
	})

	return router
}

func TestMiddlewareAuth(t *testing.T) {
	store := fakeUserStore{user: database.User{
		ID:       uuid.New(),
		Fullname: "Matteo Campana",
		ApiKey:   "a-valid-key",
	}}

	tests := []struct {
		name        string
		header      string
		store       UserStore
		wantStatus  int
		wantReached bool
	}{
		{name: "no header", store: store, wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", header: "Bearer a-valid-key", store: store, wantStatus: http.StatusUnauthorized},
		{name: "unknown key", header: "ApiKey nope", store: store, wantStatus: http.StatusUnauthorized},
		{
			name:       "database failure",
			header:     "ApiKey a-valid-key",
			store:      fakeUserStore{err: errors.New("connection refused")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:        "valid key",
			header:      "ApiKey a-valid-key",
			store:       store,
			wantStatus:  http.StatusOK,
			wantReached: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reached := false
			router := newAuthenticatedRouter(tt.store, &reached)

			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				request.Header.Set("Authorization", tt.header)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			if got := recorder.Code; got != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", got, tt.wantStatus, recorder.Body.String())
			}

			if reached != tt.wantReached {
				t.Errorf("handler reached = %t, want %t", reached, tt.wantReached)
			}
		})
	}
}

func TestUserFromContextWithoutUser(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	if _, ok := UserFromContext(c); ok {
		t.Error("UserFromContext on an unauthenticated context returned a user")
	}
}

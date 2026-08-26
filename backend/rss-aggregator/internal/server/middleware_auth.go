package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"rss-aggregator/internal/auth"
	"rss-aggregator/internal/database"

	"github.com/gin-gonic/gin"
)

// userContextKey is the key the authenticated user is stored under.
const userContextKey = "user"

// UserStore is the slice of the database layer the authentication needs.
type UserStore interface {
	GetUserByApiKey(ctx context.Context, apiKey string) (database.User, error)
}

// MiddlewareAuth rejects requests without a valid `Authorization: ApiKey <key>`
// header and stores the matching user in the request context.
func MiddlewareAuth(store UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := auth.GetAPIKey(c.Request.Header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := store.GetUserByApiKey(c, apiKey)

		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		if err != nil {
			log.Printf("auth: looking up the user by API key: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot authenticate the request"})
			return
		}

		c.Set(userContextKey, user)
		c.Next()
	}
}

// UserFromContext returns the user authenticated by MiddlewareAuth.
func UserFromContext(c *gin.Context) (database.User, bool) {
	value, exists := c.Get(userContextKey)
	if !exists {
		return database.User{}, false
	}

	user, ok := value.(database.User)

	return user, ok
}

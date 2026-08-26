package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"rss-aggregator/internal/database"
	"rss-aggregator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Every handler here derives the user from the authenticated request instead of
// reading it from the path. A follow that belongs to somebody else answers 404
// rather than 403: the caller has no business learning that it exists.

func (apiCfg *ApiConfig) CreateFeedFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		type parameters struct {
			FeedID uuid.UUID `json:"feed_id"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err := decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.FeedID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is required"})
			return
		}

		feedFollow, err := apiCfg.queries.CreateFeedFollow(c, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedID,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
	}
}

// GetFeedsFollowsHandler lists the follows of the authenticated user.
func (apiCfg *ApiConfig) GetFeedsFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		feedsFollows, err := apiCfg.queries.GetFeedsFollowsByUserId(c, user.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedsFollows))
	}
}

func (apiCfg *ApiConfig) GetFeedFollowsByIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		feedFollowID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollow, err := apiCfg.queries.GetFeedFollowByIdAndUserId(c, database.GetFeedFollowByIdAndUserIdParams{
			ID:     feedFollowID,
			UserID: user.ID,
		})

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed follow not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollow))
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		feedFollowID, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		// Look it up first: the delete reports no row count, and a follow that
		// is not the caller's must read as missing.
		_, err = apiCfg.queries.GetFeedFollowByIdAndUserId(c, database.GetFeedFollowByIdAndUserIdParams{
			ID:     feedFollowID,
			UserID: user.ID,
		})

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed follow not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollowByIdAndUserId(c, database.DeleteFeedFollowByIdAndUserIdParams{
			ID:     feedFollowID,
			UserID: user.ID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (apiCfg *ApiConfig) GetFeedFollowsByFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		feedID, err := uuid.Parse(c.Param("feed_id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is not a valid UUID"})
			return
		}

		feedFollow, err := apiCfg.queries.GetFeedFollowsByUserIdAndFeedId(c, database.GetFeedFollowsByUserIdAndFeedIdParams{
			UserID: user.ID,
			FeedID: feedID,
		})

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed follow not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollow))
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsByFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		feedID, err := uuid.Parse(c.Param("feed_id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is not a valid UUID"})
			return
		}

		_, err = apiCfg.queries.GetFeedFollowsByUserIdAndFeedId(c, database.GetFeedFollowsByUserIdAndFeedIdParams{
			UserID: user.ID,
			FeedID: feedID,
		})

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "feed follow not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollowsByUserIdAndFeedId(c, database.DeleteFeedFollowsByUserIdAndFeedIdParams{
			UserID: user.ID,
			FeedID: feedID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

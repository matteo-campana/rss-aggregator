package server

import (
	"encoding/json"
	"net/http"
	"rss-aggregator/internal/database"
	"time"

	"rss-aggregator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateFeedFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// parse and check the request parameters to create a new feed follows

		type parameters struct {
			UserID uuid.UUID `json:"user_id"`
			FeedID uuid.UUID `json:"feed_id"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err := decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.UserID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		if params.FeedID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is required"})
			return
		}

		feedFollow, err := apiCfg.queries.CreateFeedFollow(c, database.CreateFeedFollowParams{
			UserID:    params.UserID,
			FeedID:    params.FeedID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ID:        uuid.New(),
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
	}
}

func (apiCfg *ApiConfig) GetFeedsFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get all feed follows

		feedsFollows, err := apiCfg.queries.GetFeedsFollows(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedsFollows))
	}
}

func (apiCfg *ApiConfig) GetFeedsFollowsByFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get a feed follows by id

		feedFollowID := c.Param("id")

		if feedFollowID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		feedFollowID_UUID, err := uuid.Parse(feedFollowID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedsFollowsByFeedId(c, feedFollowID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
	}
}

func (apiCfg *ApiConfig) GetFeedsFollowsByUserIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get a feed follows by id

		userID := c.Param("id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		userID_UUID, err := uuid.Parse(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedsFollowsByUserId(c, userID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
	}
}

func (apiCfg *ApiConfig) GetFeedFollowsByUserIdAndFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get a feed follows by user id and feed id

		userID := c.Param("user_id")
		feedID := c.Param("feed_id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		if feedID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is required"})
			return
		}

		userID_UUID, err := uuid.Parse(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is not a valid UUID"})
			return
		}

		feedID_UUID, err := uuid.Parse(feedID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedFollowsByUserIdAndFeedId(c, database.GetFeedFollowsByUserIdAndFeedIdParams{
			UserID: userID_UUID,
			FeedID: feedID_UUID,
		})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollows))
	}
}

func (apiCfg *ApiConfig) GetFeedsFollowsByIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get a feed follows by id

		feedFollowID := c.Param("id")

		if feedFollowID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		feedFollowID_UUID, err := uuid.Parse(feedFollowID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedFollowsById(c, feedFollowID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedFollowToFeedFollow(feedFollows))
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// delete a feed follows by id

		feedFollowID := c.Param("id")

		if feedFollowID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		feedFollowID_UUID, err := uuid.Parse(feedFollowID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedFollowsById(c, feedFollowID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollows(c, feedFollowID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"message":      "Feed Follows deleted",
			"feed_follows": models.DatabaseFeedFollowToFeedFollow(feedFollows),
		})
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsByUserIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// delete a feed follows by user id

		userID := c.Param("id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		userID_UUID, err := uuid.Parse(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedsFollowsByUserId(c, userID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollowsByUserId(c, userID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"message":      "Feed Follows deleted",
			"feed_follows": models.DatabaseFeedFollowsToFeedFollows(feedFollows),
		})
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsByFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// delete a feed follows by feed id

		feedID := c.Param("id")

		if feedID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		feedID_UUID, err := uuid.Parse(feedID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedsFollowsByFeedId(c, feedID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollowsByFeedId(c, feedID_UUID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"message":      "Feed Follows deleted",
			"feed_follows": models.DatabaseFeedFollowsToFeedFollows(feedFollows),
		})
	}
}

func (apiCfg *ApiConfig) DeleteFeedFollowsByUserIdAndFeedIdHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// delete a feed follows by user id and feed id

		userID := c.Param("user_id")
		feedID := c.Param("feed_id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		if feedID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is required"})
			return
		}

		userID_UUID, err := uuid.Parse(userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is not a valid UUID"})
			return
		}

		feedID_UUID, err := uuid.Parse(feedID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "feed_id is not a valid UUID"})
			return
		}

		feedFollows, err := apiCfg.queries.GetFeedFollowsByUserIdAndFeedId(c, database.GetFeedFollowsByUserIdAndFeedIdParams{
			UserID: userID_UUID,
			FeedID: feedID_UUID,
		})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeedFollowsByUserIdAndFeedId(c, database.DeleteFeedFollowsByUserIdAndFeedIdParams{
			UserID: userID_UUID,
			FeedID: feedID_UUID,
		})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"message":      "Feed Follows deleted",
			"feed_follows": models.DatabaseFeedFollowToFeedFollow(feedFollows),
		})
	}
}

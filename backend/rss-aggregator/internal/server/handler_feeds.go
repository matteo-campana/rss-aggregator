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

func (apiCfg *ApiConfig) CreateFeedHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// parse and check the request parameters to create a new feed

		type parameters struct {
			Url string `json:"url"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err := decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.Url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		feed, err := apiCfg.queries.CreateFeed(c, database.CreateFeedParams{
			Url:       params.Url,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ID:        uuid.New(),
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.DatabaseFeedToFeed(feed))
	}
}

func (apiCfg *ApiConfig) GetFeedsHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get all feeds

		feeds, err := apiCfg.queries.GetFeeds(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
	}
}

func (apiCfg *ApiConfig) GetFeedHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// get a feed by id

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

		feed, err := apiCfg.queries.GetFeed(c, feedID_UUID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedToFeed(feed))
	}
}

func (apiCfg *ApiConfig) UpdateFeedHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// update a feed by id

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

		type parameters struct {
			Url string `json:"url"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err = decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.Url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		feed, err := apiCfg.queries.UpdateFeed(c, database.UpdateFeedParams{
			ID:        feedID_UUID,
			Url:       params.Url,
			UpdatedAt: time.Now().UTC(),
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseFeedToFeed(feed))
	}
}

func (apiCfg *ApiConfig) DeleteFeedHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// delete a feed by id

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

		feed, err := apiCfg.queries.GetFeed(c, feedID_UUID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = apiCfg.queries.DeleteFeed(c, feedID_UUID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"message": "feed deleted",
			"feed":    models.DatabaseFeedToFeed(feed),
		})
	}
}

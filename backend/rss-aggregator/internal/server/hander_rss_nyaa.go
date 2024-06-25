package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rss-aggregator/internal/models"
	rss "rss-aggregator/pkg/rss"
)

func (apiCfg *ApiConfig) HandlerNyaaRss() gin.HandlerFunc {
	return func(c *gin.Context) {

		// language Language, resolution Resolution

		fetchAndParseRSSRequest := &rss.FetchAndParseRSSRequest{
			Language:   rss.ENG,
			Resolution: rss.Resolution1080p,
		}

		db_channel, items, err := rss.FetchAndParseRSS(c, fetchAndParseRSSRequest, apiCfg.queries)

		if err != nil {
			c.JSON(501, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"channel": db_channel,
			"items":   models.DatabaseItemsToItems(*items),
		})
	}
}

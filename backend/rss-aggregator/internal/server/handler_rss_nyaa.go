package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rss-aggregator/internal/models"
	rss "rss-aggregator/pkg/rss"
)

// HandlerNyaaRss fetches a nyaa.si feed and stores it. Language and resolution
// come from the query string; the previous defaults apply when they are absent.
func (apiCfg *ApiConfig) HandlerNyaaRss() gin.HandlerFunc {
	return func(c *gin.Context) {

		language, err := rss.ParseLanguage(c.Query("language"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resolution, err := rss.ParseResolution(c.Query("resolution"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fetchAndParseRSSRequest := &rss.FetchAndParseRSSRequest{
			Language:   language,
			Resolution: resolution,
		}

		db_channel, items, err := rss.FetchAndParseRSS(c, fetchAndParseRSSRequest, apiCfg.queries)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"channel": models.DatabaseChannelToChannel(db_channel),
			"items":   models.DatabaseItemsToItems(items),
		})
	}
}

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// itemsNamespace groups everything derived from the scraped items — the item
// listing, its categories and the channels. The scraper bumps this one version
// after a run, which retires all three at once.
const itemsNamespace = "items"

const jsonContentType = "application/json; charset=utf-8"

// serveCached answers from Redis when it can, and otherwise runs produce,
// caches its JSON and returns it. Without Redis it is just produce.
//
// Validation belongs in the caller: only the database work goes through here,
// so a rejected request never reaches the cache.
func (apiCfg *ApiConfig) serveCached(c *gin.Context, endpoint string, produce func() (any, error)) {
	key := apiCfg.Cache().Key(c, itemsNamespace, endpoint, c.Request.URL.RawQuery)

	if payload, ok := apiCfg.Cache().GetRaw(c, key); ok {
		c.Header("X-Cache", "hit")
		c.Data(http.StatusOK, jsonContentType, payload)

		return
	}

	value, err := produce()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payload, err := json.Marshal(value)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiCfg.Cache().SetRaw(c, key, payload)

	c.Header("X-Cache", "miss")
	c.Data(http.StatusOK, jsonContentType, payload)
}

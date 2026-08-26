package server

import (
	"net/http"
	"strconv"
	"time"

	"rss-aggregator/internal/auth"

	"github.com/gin-gonic/gin"
)

// MiddlewareRateLimit counts requests per API key. It runs after
// MiddlewareAuth, so an unauthenticated request is already rejected and cannot
// consume anybody's allowance.
//
// Without Redis every request is allowed: a limiter that cannot count must not
// lock everybody out.
func (apiCfg *ApiConfig) MiddlewareRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, err := auth.GetAPIKey(c.Request.Header)

		if err != nil {
			// Unauthenticated requests never get here; if one does, let the
			// auth middleware be the one to refuse it.
			c.Next()
			return
		}

		decision := apiCfg.cache.Allow(c, apiKey)

		c.Header("X-RateLimit-Limit", strconv.Itoa(decision.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(decision.Remaining))

		if decision.Allowed {
			c.Next()
			return
		}

		retryAfter := int(decision.RetryAfter / time.Second)

		if retryAfter < 1 {
			retryAfter = 1
		}

		c.Header("Retry-After", strconv.Itoa(retryAfter))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "rate limit exceeded, retry in " + strconv.Itoa(retryAfter) + "s",
		})
	}
}

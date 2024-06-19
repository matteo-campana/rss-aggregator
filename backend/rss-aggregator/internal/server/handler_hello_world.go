package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (apiCfg *ApiConfig) HelloWorldHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := make(map[string]string)
		resp["message"] = "Hello World"
		c.JSON(http.StatusOK, resp)
	}
}

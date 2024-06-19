package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (s *Server) RegisterRoutes() http.Handler {
func (apiCfg *ApiConfig) RegisterRoutes() http.Handler {
	router := gin.Default()

	routerGroup := router.Group("/api/v1")

	// health

	routerGroup.GET("/", apiCfg.HelloWorldHandler())
	routerGroup.GET("/health", apiCfg.HealthHandler())

	// users
	routerGroupUsers := routerGroup.Group("/users")

	routerGroupUsers.POST("/", apiCfg.CreateUserHandler())
	routerGroupUsers.GET("/:id", apiCfg.GetUserHandler())
	routerGroupUsers.GET("/", apiCfg.GetUsersHandler())
	routerGroupUsers.PUT("/:id", apiCfg.UpdateUserHandler())
	routerGroupUsers.DELETE("/:id", apiCfg.DeleteUserHandler())

	// feeds

	routerGroupFeeds := routerGroup.Group("/feeds")
	routerGroupFeeds.POST("/", apiCfg.CreateFeedHandler())
	routerGroupFeeds.GET("/:id", apiCfg.GetFeedHandler())
	routerGroupFeeds.GET("/", apiCfg.GetFeedsHandler())
	routerGroupFeeds.PUT("/:id", apiCfg.UpdateFeedHandler())
	routerGroupFeeds.DELETE("/:id", apiCfg.DeleteFeedHandler())

	// feed follows

	routerGroupFeedFollows := routerGroup.Group("/feed-follows")

	routerGroupFeedFollows.POST("/", apiCfg.CreateFeedFollowsHandler())
	routerGroupFeedFollows.GET("/", apiCfg.GetFeedsFollowsHandler())
	routerGroupFeedFollows.DELETE("/:id", apiCfg.DeleteFeedFollowsHandler())

	routerGroupFeedFollows.GET("/user/:user_id", apiCfg.GetFeedsFollowsByUserIdHandler())
	routerGroupFeedFollows.GET("/feed/:feed_id", apiCfg.GetFeedsFollowsByFeedIdHandler())

	routerGroupFeedFollows.GET("/user/:user_id/feed/:feed_id", apiCfg.GetFeedFollowsByUserIdAndFeedIdHandler())
	routerGroupFeedFollows.DELETE("/user/:user_id/feed/:feed_id", apiCfg.DeleteFeedFollowsByUserIdAndFeedIdHandler())

	return router
}

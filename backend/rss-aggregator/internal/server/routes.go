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

	// registration is the only public users route: it is how a caller gets
	// the API key every other route requires.

	routerGroup.POST("/users/", apiCfg.CreateUserHandler())

	// everything below requires `Authorization: ApiKey <key>`

	authenticated := routerGroup.Group("")
	authenticated.Use(MiddlewareAuth(apiCfg.queries))

	// users
	routerGroupUsers := authenticated.Group("/users")

	routerGroupUsers.GET("/me", apiCfg.GetCurrentUserHandler())
	routerGroupUsers.GET("/:id", apiCfg.GetUserHandler())
	routerGroupUsers.GET("/", apiCfg.GetUsersHandler())
	routerGroupUsers.PUT("/:id", apiCfg.UpdateUserHandler())
	routerGroupUsers.DELETE("/:id", apiCfg.DeleteUserHandler())

	// feeds

	routerGroupFeeds := authenticated.Group("/feeds")
	routerGroupFeeds.POST("/", apiCfg.CreateFeedHandler())
	routerGroupFeeds.GET("/:id", apiCfg.GetFeedHandler())
	routerGroupFeeds.GET("/", apiCfg.GetFeedsHandler())
	routerGroupFeeds.PUT("/:id", apiCfg.UpdateFeedHandler())
	routerGroupFeeds.DELETE("/:id", apiCfg.DeleteFeedHandler())

	// feed follows

	routerGroupFeedFollows := authenticated.Group("/feed-follows")

	routerGroupFeedFollows.POST("/", apiCfg.CreateFeedFollowsHandler())
	routerGroupFeedFollows.GET("/", apiCfg.GetFeedsFollowsHandler())
	routerGroupFeedFollows.DELETE("/:id", apiCfg.DeleteFeedFollowsHandler())

	routerGroupFeedFollows.GET("/user/:user_id", apiCfg.GetFeedsFollowsByUserIdHandler())
	routerGroupFeedFollows.GET("/feed/:feed_id", apiCfg.GetFeedsFollowsByFeedIdHandler())

	routerGroupFeedFollows.GET("/user/:user_id/feed/:feed_id", apiCfg.GetFeedFollowsByUserIdAndFeedIdHandler())
	routerGroupFeedFollows.DELETE("/user/:user_id/feed/:feed_id", apiCfg.DeleteFeedFollowsByUserIdAndFeedIdHandler())

	// items and channels: what the scraper has already stored

	routerGroupItems := authenticated.Group("/items")

	routerGroupItems.GET("/", apiCfg.ListItemsHandler())
	routerGroupItems.GET("/categories", apiCfg.ListItemCategoriesHandler())

	authenticated.GET("/channels/", apiCfg.GetChannelsHandler())

	// nyaa

	routerGroupNyaa := authenticated.Group("/nyaa")
	routerGroupNyaa.GET("/rss", apiCfg.HandlerNyaaRss())

	return router
}

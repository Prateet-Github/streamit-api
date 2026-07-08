package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterSubscriptionRoutes(
	router *gin.Engine,
	subscriptionHandler *handlers.SubscriptionHandler,
	jwtSecret string,
) {
	subs := router.Group("/api/subscriptions")

	subs.Use(middlewares.OptionalAuth(jwtSecret))

	subs.GET("/:channelId", subscriptionHandler.GetSubscriptionStatus)

	protected := router.Group("/api/subscriptions")

	protected.Use(middlewares.Auth(jwtSecret))

	protected.GET("/", subscriptionHandler.GetMySubscriptions)
	protected.POST("/:channelId", subscriptionHandler.Subscribe)
	protected.DELETE("/:channelId", subscriptionHandler.Unsubscribe)

}

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

	subs.Use(middlewares.Auth(jwtSecret))

	subs.GET("/", subscriptionHandler.GetMySubscriptions)
	subs.GET("/:channelId", subscriptionHandler.GetSubscriptionStatus)
	subs.POST("/:channelId", subscriptionHandler.Subscribe)
	subs.DELETE("/:channelId", subscriptionHandler.Unsubscribe)
}

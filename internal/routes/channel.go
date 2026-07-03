package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterChannelRoutes(
	router *gin.Engine,
	channelHandler *handlers.ChannelHandler,
) {
	channels := router.Group("/api/channels")

	channels.GET("/:username", channelHandler.GetChannel)
}

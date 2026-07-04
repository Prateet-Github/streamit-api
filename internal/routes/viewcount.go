package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterViewCountRoutes(
	router *gin.Engine,
	viewCountHandler *handlers.ViewCountHandler,
	jwtSecret string,
) {
	viewcount := router.Group("/api/videocount")

	viewcount.Use(middlewares.Auth(jwtSecret))

	viewcount.POST("/:id/view", viewCountHandler.Heartbeat)
}

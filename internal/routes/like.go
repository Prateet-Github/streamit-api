package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterLikeRoutes(
	router *gin.Engine,
	likeHandler *handlers.LikeHandler,
	jwtSecret string,
) {
	likes := router.Group("/api/videos")
	likes.Use(middlewares.Auth(jwtSecret))

	likes.POST("/:videoId/like", likeHandler.Like)
	likes.DELETE("/:videoId/like", likeHandler.Unlike)
}

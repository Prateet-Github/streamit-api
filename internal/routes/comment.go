package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(
	router *gin.Engine,
	commentHandler *handlers.CommentHandler,
	jwtSecret string,
) {
	auth := router.Group("/api/comments")
	auth.Use(middlewares.Auth(jwtSecret))

	auth.POST("/video/:videoId", commentHandler.CreateComment)
	auth.POST("/:commentId/replies", commentHandler.CreateReply)
	auth.DELETE("/:commentId", commentHandler.DeleteComment)

	router.GET("/api/comments/video/:videoId", commentHandler.GetComments)
	router.GET("/api/comments/:commentId/replies", commentHandler.GetReplies)
}

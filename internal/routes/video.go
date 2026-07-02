package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterVideoRoutes(
	router *gin.Engine,
	videoHandler *handlers.VideoHandler,
	jwtSecret string,
) {
	videos := router.Group("/api/video")

	videos.GET("/", videoHandler.GetAllVideos)
	videos.GET("/:id", videoHandler.GetVideoByID)

	videos.Use(middlewares.Auth(jwtSecret))

	videos.POST("/upload-url", videoHandler.GetUploadURL)
	videos.POST("/confirm-upload", videoHandler.ConfirmUpload)

	internal := router.Group("/internal/videos")

	internal.POST("/:id/complete", videoHandler.CompleteVideo)
}

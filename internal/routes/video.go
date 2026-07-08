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
	videos.GET("/search", videoHandler.SearchVideos)
	videos.GET("/:id", videoHandler.GetVideoByID)

	auth := videos.Group("")
	auth.Use(middlewares.Auth(jwtSecret))

	auth.GET("/my-videos", videoHandler.GetMyVideos)
	auth.POST("/upload-url", videoHandler.GetUploadURL)
	auth.POST("/confirm-upload", videoHandler.ConfirmUpload)

	internal := router.Group("/internal/videos")

	internal.POST("/:id/complete", videoHandler.CompleteVideo)
}

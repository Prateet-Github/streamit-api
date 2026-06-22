package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", handlers.Health)
}

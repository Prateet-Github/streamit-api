package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler) {
	auth := router.Group("/api/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}

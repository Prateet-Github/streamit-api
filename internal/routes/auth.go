package routes

import (
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	jwtSecret string,
) {
	auth := router.Group("/api/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := auth.Group("")
	protected.Use(middlewares.Auth(jwtSecret))

	protected.GET("/me", authHandler.Me)

}

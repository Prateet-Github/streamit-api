package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/routes"
)

func New(cfg *config.Config) *gin.Engine {
	db, err := database.Connect(
		cfg.MongoURI,
		cfg.DatabaseName,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB connected")

	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)

	router := gin.Default()

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router, authHandler, cfg.JWTSecret)

	return router
}

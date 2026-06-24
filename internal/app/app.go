package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/routes"
	"github.com/Prateet-Github/streamit-api/internal/s3"
)

func New(cfg *config.Config) *gin.Engine {
	db, err := database.Connect(
		cfg.MongoURI,
		cfg.DatabaseName,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := database.CreateIndexes(db); err != nil {
		log.Fatal(err)
	}

	s3Client, err := s3.NewClient(cfg)

	if err != nil {
		log.Fatal(err)
	}

	_ = s3Client

	log.Println("S3 client connected")
	log.Println("MongoDB connected")
	log.Println("Indexes created")

	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)

	videoRepo := repositories.NewVideoRepository(db)
	videoHandler := handlers.NewVideoHandler(s3Client, cfg, videoRepo)

	router := gin.Default()

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router, authHandler, cfg.JWTSecret)
	routes.RegisterVideoRoutes(router, videoHandler, cfg.JWTSecret)

	return router
}

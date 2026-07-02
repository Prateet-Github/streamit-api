package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/queue"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/routes"
	"github.com/Prateet-Github/streamit-api/internal/s3"
)

func New(cfg *config.Config) *gin.Engine {
	// 1. MongoDB
	db, err := database.Connect(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}
	log.Println("MongoDB connected successfully")

	if err := database.CreateIndexes(db); err != nil {
		log.Fatal("MongoDB index creation failed: ", err)
	}
	log.Println("MongoDB database indexes created")

	// 2. AWS S3
	s3Client, err := s3.NewClient(cfg)
	if err != nil {
		log.Fatal("S3 client initialization failed: ", err)
	}
	log.Println("S3 client connected successfully")

	// 3. Redis Core
	redisClient := queue.NewRedisClient(cfg)
	if err := queue.Ping(redisClient); err != nil {
		log.Fatal("Redis ping failed: ", err)
	}
	log.Println("Redis client connected and verified")
	redisClient.Close()

	// 4. Asynq Client
	asynqClient := queue.NewAsynqClient(cfg)
	log.Println("Asynq client initialized successfully")

	// 5. Wire Repositories & Handlers
	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)

	videoRepo := repositories.NewVideoRepository(db)
	videoHandler := handlers.NewVideoHandler(s3Client, cfg, videoRepo, asynqClient)

	likeRepo := repositories.NewLikeRepository(db.DB)
	likeHandler := handlers.NewLikeHandler(likeRepo, videoRepo)

	// 6. Router Setup
	router := gin.Default()

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router, authHandler, cfg.JWTSecret)
	routes.RegisterVideoRoutes(router, videoHandler, cfg.JWTSecret)
	routes.RegisterLikeRoutes(router, likeHandler, cfg.JWTSecret)

	return router
}

package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/handlers"
	"github.com/Prateet-Github/streamit-api/internal/queue"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/routes"
	"github.com/Prateet-Github/streamit-api/internal/s3"
	"github.com/Prateet-Github/streamit-api/internal/services/viewcount"
	"github.com/gin-contrib/cors"
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
	// redisClient.Close()

	// 4. Asynq Client
	asynqClient := queue.NewAsynqClient(cfg)
	log.Println("Asynq client initialized successfully")

	// 5. Repositories & Handlers
	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)

	videoRepo := repositories.NewVideoRepository(db)
	videoHandler := handlers.NewVideoHandler(s3Client, cfg, videoRepo, userRepo, asynqClient)

	likeRepo := repositories.NewLikeRepository(db.DB)
	likeHandler := handlers.NewLikeHandler(likeRepo, videoRepo)

	commentRepo := repositories.NewCommentRepository(db.DB)
	commentHandler := handlers.NewCommentHandler(commentRepo, videoRepo, userRepo)

	subscriptionRepo := repositories.NewSubscriptionRepository(db.DB)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionRepo, userRepo)

	channelHandler := handlers.NewChannelHandler(userRepo, videoRepo)

	producer := viewcount.NewProducer(redisClient)
	viewCountHandler := handlers.NewViewCountHandler(producer)

	validator := viewcount.NewValidator(redisClient)
	deduplicator := viewcount.NewDeduplicator(redisClient)
	analytics := viewcount.NewAnalytics(redisClient)
	counter := viewcount.NewCounter(redisClient)

	worker := viewcount.NewWorker(redisClient, validator, deduplicator, analytics, counter)

	ctx := context.Background()

	if err := worker.CreateGroup(ctx); err != nil {
		log.Fatal(err)
	}

	go worker.Start(ctx)

	viewRepository := repositories.NewViewRepository(db)

	flusher := viewcount.NewFlusher(
		redisClient,
		viewRepository,
	)

	cron := viewcount.NewCron(flusher)

	go cron.Start(ctx)

	// 6. Router Setup
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router, authHandler, cfg.JWTSecret)
	routes.RegisterVideoRoutes(router, videoHandler, cfg.JWTSecret)
	routes.RegisterLikeRoutes(router, likeHandler, cfg.JWTSecret)
	routes.RegisterCommentRoutes(router, commentHandler, cfg.JWTSecret)
	routes.RegisterSubscriptionRoutes(router, subscriptionHandler, cfg.JWTSecret)
	routes.RegisterViewCountRoutes(router, viewCountHandler, cfg.JWTSecret)
	routes.RegisterChannelRoutes(router, channelHandler)
	return router
}

package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/routes"
)

func New(cfg *config.Config) *gin.Engine {
	client, err := database.Connect(cfg.MongoURI)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB connected")

	_ = client

	router := gin.Default()

	routes.RegisterHealthRoutes(router)

	return router
}

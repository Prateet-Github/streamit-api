package main

import (
	"log"

	"github.com/Prateet-Github/streamit-api/internal/app"
	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(gin.ReleaseMode) // Release mode for production

	server := app.New(cfg)

	log.Printf("Server is running on port: %s (Release Mode)", cfg.Port)

	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

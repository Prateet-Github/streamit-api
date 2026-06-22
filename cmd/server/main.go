package main

import (
	"log"

	"github.com/Prateet-Github/streamit-api/internal/app"
	"github.com/Prateet-Github/streamit-api/internal/config"
)

func main() {
	cfg := config.Load()

	server := app.New()

	log.Printf("Server running on :%s", cfg.Port)

	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

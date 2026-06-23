package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	MongoURI     string
	DatabaseName string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found, using system env")
	}

	return &Config{
		Port:         os.Getenv("PORT"),
		Env:          os.Getenv("APP_ENV"),
		MongoURI:     os.Getenv("MONGODB_URI"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
	}
}

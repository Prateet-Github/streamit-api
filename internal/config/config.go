package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Env                string
	MongoURI           string
	DatabaseName       string
	JWTSecret          string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	S3RawBucket  string
	S3ProdBucket string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	FrontendURL string
}

func Load() *Config {
	_ = godotenv.Load() // for prod

	// if err != nil {
	// 	log.Println(".env not found, using system env")
	// }

	return &Config{
		Port:               os.Getenv("PORT"),
		Env:                os.Getenv("APP_ENV"),
		MongoURI:           os.Getenv("MONGODB_URI"),
		DatabaseName:       os.Getenv("DATABASE_NAME"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3RawBucket:        os.Getenv("S3_RAW_BUCKET"),
		S3ProdBucket:       os.Getenv("S3_PROD_BUCKET"),
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		FrontendURL:        os.Getenv("FRONTEND_URL"),
	}
}

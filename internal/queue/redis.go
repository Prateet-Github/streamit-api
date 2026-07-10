package queue

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
}

func Ping(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}

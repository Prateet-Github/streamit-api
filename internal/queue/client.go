package queue

import (
	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/hibiken/asynq"
)

func NewAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
			Password: cfg.RedisPassword,
		},
	)
}

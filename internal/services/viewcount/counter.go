package viewcount

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Counter struct {
	redis *redis.Client
}

func NewCounter(redisClient *redis.Client) *Counter {
	return &Counter{
		redis: redisClient,
	}
}

func (c *Counter) Increment(
	ctx context.Context,
	event ViewEvent,
) error {

	return c.redis.HIncrBy(
		ctx,
		"streamit:views:hot_counters",
		event.VideoID,
		1,
	).Err()
}

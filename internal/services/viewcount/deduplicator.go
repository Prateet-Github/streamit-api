package viewcount

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Deduplicator struct {
	redis *redis.Client
}

func NewDeduplicator(redisClient *redis.Client) *Deduplicator {
	return &Deduplicator{
		redis: redisClient,
	}
}

func (d *Deduplicator) Check(
	ctx context.Context,
	event ViewEvent,
) (bool, error) {

	key := "view:" + event.VideoID + ":" + event.ViewerID

	ok, err := d.redis.SetNX(
		ctx,
		key,
		1,
		4*time.Hour, // TTL of 4 hours for the view event
	).Result()

	if err != nil {
		return false, err
	}

	return ok, nil
}

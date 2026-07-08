package viewcount

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Analytics struct {
	redis *redis.Client
}

func NewAnalytics(redisClient *redis.Client) *Analytics {
	return &Analytics{
		redis: redisClient,
	}
}

func (a *Analytics) Record(
	ctx context.Context,
	event ViewEvent,
) error {

	date := time.Now().Format("2006-01-02")

	key := "unique:" + date + ":" + event.VideoID

	if err := a.redis.PFAdd(
		ctx,
		key,
		event.ViewerID,
	).Err(); err != nil {
		return err
	}

	return a.redis.Expire(
		ctx,
		key,
		48*time.Hour,
	).Err()
}

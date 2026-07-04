package viewcount

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ViewEvent struct {
	VideoID   string
	ViewerID  string
	Elapsed   int
	Timestamp time.Time
}

type Producer struct {
	redis *redis.Client
}

func NewProducer(redisClient *redis.Client) *Producer {
	return &Producer{
		redis: redisClient,
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	event ViewEvent,
) error {

	err := p.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "streamit:view_events",
		Values: map[string]any{
			"videoId":   event.VideoID,
			"viewerId":  event.ViewerID,
			"elapsed":   event.Elapsed,
			"timestamp": event.Timestamp.Unix(),
		},
	}).Err()

	return err

}

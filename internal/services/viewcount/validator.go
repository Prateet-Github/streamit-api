package viewcount

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Validator struct {
	redis *redis.Client
}

func NewValidator(redisClient *redis.Client) *Validator {
	return &Validator{
		redis: redisClient,
	}
}

func (v *Validator) Validate(
	ctx context.Context,
	event ViewEvent,
) (bool, error) {

	key := "track:" + event.ViewerID + ":" + event.VideoID

	if err := v.redis.SAdd(ctx, key, event.Elapsed).Err(); err != nil {
		return false, err
	}

	if err := v.redis.Expire(ctx, key, 5*time.Minute).Err(); err != nil { // TTL of 5 minutes for the set
		return false, err
	}

	count, err := v.redis.SCard(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count == 3, nil
}

func (v *Validator) Cleanup(
	ctx context.Context,
	event ViewEvent,
) error {

	key := "track:" + event.ViewerID + ":" + event.VideoID

	return v.redis.Del(ctx, key).Err()
}

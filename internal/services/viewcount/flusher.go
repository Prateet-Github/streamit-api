package viewcount

import (
	"context"
	"strconv"

	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type Flusher struct {
	redis      *redis.Client
	repository *repositories.ViewRepository
}

func NewFlusher(
	redisClient *redis.Client,
	repository *repositories.ViewRepository,
) *Flusher {
	return &Flusher{
		redis:      redisClient,
		repository: repository,
	}
}

func (f *Flusher) Flush(ctx context.Context) error {

	counts, err := f.redis.HGetAll(
		ctx,
		"streamit:views:hot_counters",
	).Result()
	if err != nil {
		return err
	}

	if len(counts) == 0 {
		return nil
	}

	increments := make(map[string]int64)

	for videoID, countStr := range counts {
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			return err
		}

		increments[videoID] = count
	}

	if err := f.repository.BulkIncrementViews(ctx, increments); err != nil {
		return err
	}

	for videoID, count := range increments {
		if err := f.redis.HIncrBy(
			ctx,
			"streamit:views:hot_counters",
			videoID,
			-count,
		).Err(); err != nil {
			return err
		}
	}

	return nil
}

package viewcount

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	streamName   = "streamit:view_events"
	groupName    = "viewcount-group"
	consumerName = "worker-1"
)

type Worker struct {
	redis *redis.Client
}

func NewWorker(redisClient *redis.Client) *Worker {
	return &Worker{
		redis: redisClient,
	}
}

func (w *Worker) CreateGroup(ctx context.Context) error {
	err := w.redis.XGroupCreateMkStream(
		ctx,
		streamName,
		groupName,
		"$",
	).Err()

	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}

	return nil
}

func (w *Worker) Start(ctx context.Context) {
	for {
		streams, err := w.redis.XReadGroup(
			ctx,
			&redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			},
		).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}

			fmt.Println("worker error:", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {

				fmt.Println("Event ID:", msg.ID)
				fmt.Println("Values:", msg.Values)

				if err := w.redis.XAck(
					ctx,
					streamName,
					groupName,
					msg.ID,
				).Err(); err != nil {
					fmt.Println("ack error:", err)
				}
			}
		}
	}
}
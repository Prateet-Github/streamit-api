package viewcount

import (
	"context"
	"fmt"
	"strconv"
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
	redis     *redis.Client
	validator *Validator
}

func NewWorker(redisClient *redis.Client, validator *Validator) *Worker {
	return &Worker{
		redis:     redisClient,
		validator: validator,
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

				// fmt.Println("Values:", msg.Values)
				event, err := decodeViewEvent(msg.Values)
				if err != nil {
					fmt.Println("decode error:", err)
					continue
				}

				fmt.Printf("%+v\n", event)

				ok, err := w.validator.Validate(ctx, event)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println("Validation:", ok)

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

func decodeViewEvent(values map[string]any) (ViewEvent, error) {
	videoID, _ := values["videoId"].(string)
	viewerID, _ := values["viewerId"].(string)

	elapsedStr, _ := values["elapsed"].(string)
	elapsed, err := strconv.Atoi(elapsedStr)
	if err != nil {
		return ViewEvent{}, err
	}

	timestampStr, _ := values["timestamp"].(string)
	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return ViewEvent{}, err
	}

	return ViewEvent{
		VideoID:   videoID,
		ViewerID:  viewerID,
		Elapsed:   elapsed,
		Timestamp: time.Unix(ts, 0),
	}, nil
}

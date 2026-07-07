package viewcount

import (
	"context"
	"log"
	"time"
)

type Cron struct {
	flusher *Flusher
}

func NewCron(flusher *Flusher) *Cron {
	return &Cron{
		flusher: flusher,
	}
}

func (c *Cron) Start(ctx context.Context) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-ticker.C:

			if err := c.flusher.Flush(ctx); err != nil {
				log.Println("flush error:", err)
			}

		case <-ctx.Done():
			return
		}
	}
}

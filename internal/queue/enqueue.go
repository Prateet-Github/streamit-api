package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

func NewProcessVideoTask(payload VideoTask) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeProcessVideo, data), nil
}

package save_race_time

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/pgqueue"
)

type ProducerSender interface {
	SendMessage(ctx context.Context, key string, message interface{}) error
}

type raceSetter struct {
	kafkaProducer ProducerSender
}

func NewRaceSetter(kafkaProducer ProducerSender) raceSetter {
	return raceSetter{
		kafkaProducer: kafkaProducer,
	}
}

type Payload struct {
	UserName  string    `json:"user_name"`
	TimeStamp time.Time `json:"time_stamp"`
}

func (tts raceSetter) Handle(ctx context.Context, taskName pgqueue.TaskKind, payloadData []byte) error {
	if taskName != pgqueue.RaceTime {
		return fmt.Errorf("wrong task handler: got %s, want %s", taskName, pgqueue.RaceTime)
	}

	var payload Payload
	err := json.Unmarshal(payloadData, &payload)
	if err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}

	err = tts.kafkaProducer.SendMessage(ctx, payload.UserName, payload)
	if err != nil {
		return fmt.Errorf("cannot send race time: %w", err)
	}
	return nil
}

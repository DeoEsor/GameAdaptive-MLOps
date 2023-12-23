package test_time_setter

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

type testTimeSetter struct {
	kafkaProducer ProducerSender
}

func NewTestTimeSetter(kafkaProducer ProducerSender) testTimeSetter {
	return testTimeSetter{
		kafkaProducer: kafkaProducer,
	}
}

type Payload struct {
	Key       string
	TimeStamp time.Time
}

func (tts testTimeSetter) Handle(ctx context.Context, taskName pgqueue.TaskKind, payloadData []byte) error {
	if taskName != pgqueue.TimeSetKind {
		return fmt.Errorf("wrong task handler: got %s, want %s", taskName, pgqueue.TimeSetKind)
	}

	var payload Payload
	err := json.Unmarshal(payloadData, &payload)
	if err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}

	err = tts.kafkaProducer.SendMessage(ctx, payload.Key, payload.TimeStamp)
	if err != nil {
		return fmt.Errorf("cannot test time: %w", err)
	}
	return nil
}

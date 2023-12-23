package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app/usecase/test_time_setter"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/pgqueue"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (i Implementation) SendTestTime(ctx context.Context, req *domain.SendTestTimeRequest) (*domain.SendTestTimeResponse, error) {
	logrus.Info(req.Payload.Timestamp)

	payload := test_time_setter.Payload{
		Key:       uuid.NewString(),
		TimeStamp: req.Payload.Timestamp,
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload to pgqueue: %w", err)
	}

	err = i.scheduler.Schedule(ctx, pgqueue.TimeSetKind, payloadStr)
	if err != nil {
		return nil, fmt.Errorf("cannot send message: %w", err)
	}
	logrus.Info("schaduled message")
	return nil, nil
}

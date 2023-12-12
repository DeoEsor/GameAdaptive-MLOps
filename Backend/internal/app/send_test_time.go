package app

import (
	"context"
	"fmt"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/sirupsen/logrus"
)

func (i Implementation) SendTestTime(ctx context.Context, req *domain.SendTestTimeRequest) (*domain.SendTestTimeResponse, error) {
	logrus.Info(req.Payload.Timestamp)
	err := i.producerSender.SendMessage(ctx, req.Payload.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("cannot send message: %w", err)
	}
	return nil, nil
}

package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/domain"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app/usecase/save_race_time"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/pgqueue"
	"github.com/sirupsen/logrus"
)

func (i Implementation) SaveRaceTime(ctx context.Context, req *domain.SaveRaceRequest) (*domain.SaveRaceResponse, error) {
	logrus.Infof("user %s has time %s", req.Payload.UserName, req.Payload.Timestamp)

	payload := save_race_time.Payload{
		UserName:  req.Payload.UserName,
		TimeStamp: req.Payload.Timestamp,
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload for pgqueue: %w", err)
	}

	err = i.scheduler.Schedule(ctx, pgqueue.RaceTime, payloadStr)
	if err != nil {
		return nil, fmt.Errorf("cannot send message: %w", err)
	}
	logrus.Info("schaduled message")
	return nil, nil
}

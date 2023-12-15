package app

import (
	"context"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/pgqueue"
)

type (
	PgqueueScheduler interface {
		Schedule(ctx context.Context, taskName pgqueue.TaskKind, payload []byte) error
	}
)

// Implementation структура для реализации различных ручек
type Implementation struct {
	scheduler PgqueueScheduler
}

// NewImplementation конструктор для Implementation
func NewImplementation(
	producerSender PgqueueScheduler,
) (*Implementation, error) {
	return &Implementation{
		scheduler: producerSender,
	}, nil
}

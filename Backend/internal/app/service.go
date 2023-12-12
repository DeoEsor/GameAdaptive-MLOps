package app

import (
	"context"
	"time"
)

type (
	ProducerSender interface {
		SendMessage(ctx context.Context, message time.Time) error
	}
)

// Implementation структура для реализации различных ручек
type Implementation struct {
	producerSender ProducerSender
}

// NewImplementation конструктор для Implementation
func NewImplementation(
	producerSender ProducerSender,
) (*Implementation, error) {
	return &Implementation{
		producerSender: producerSender,
	}, nil
}

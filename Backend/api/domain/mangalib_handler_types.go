package domain

import (
	"context"
	"time"
)

type HandlerType string

const (
	SendTestTime = HandlerType("SendTestTime")
)

// Интерфейс со всеми ручками сервиса
type Handlers interface {
	// SendTestTime - отправляет время в кафку брокер
	SendTestTime(ctx context.Context, req *SendTestTimeRequest) (*SendTestTimeResponse, error)
}

type SendTestTimePayload struct {
	Timestamp time.Time `json:"time_stamp"`
}

type SendTestTimeRequest struct {
	Payload *SendTestTimePayload `in:"body=json"`
}

type SendTestTimeResponse struct {
}

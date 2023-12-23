package domain

import (
	"context"
	"time"
)

type HandlerType string

const (
	SendTestTime = HandlerType("SendTestTime")
	SaveRaceTime = HandlerType("SaveRaceTime")
)

// Интерфейс со всеми ручками сервиса
type Handlers interface {
	// SendTestTime - отправляет время в кафку брокер
	SendTestTime(ctx context.Context, req *SendTestTimeRequest) (*SendTestTimeResponse, error)
	SaveRaceTime(ctx context.Context, req *SaveRaceRequest) (*SaveRaceResponse, error)
}

type SaveRacePayload struct {
	UserName  string    `json:"user_name"`
	Timestamp time.Time `json:"time_stamp"`
}

type SaveRaceRequest struct {
	Payload *SaveRacePayload `in:"body=json"`
}

type SaveRaceResponse struct{}

type SendTestTimePayload struct {
	Timestamp time.Time `json:"time_stamp"`
}

type SendTestTimeRequest struct {
	Payload *SendTestTimePayload `in:"body=json"`
}

type SendTestTimeResponse struct{}

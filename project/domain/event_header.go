package domain

import (
	"github.com/google/uuid"
	"time"
)

type EventHeader struct {
	ID             string    `json:"id"`
	PublishedAt    time.Time `json:"published_at"`
	IdempotencyKey string    `json:"idempotency_key"`
}

func NewHeader(idempotencyKey string) EventHeader {
	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}
	return EventHeader{
		ID:             uuid.NewString(),
		PublishedAt:    time.Now().UTC(),
		IdempotencyKey: idempotencyKey,
	}
}

type Event interface {
	isEvent()
}

package domain

import (
	"github.com/google/uuid"
	"time"
)

type EventHeader struct {
	ID          string    `json:"id"`
	PublishedAt time.Time `json:"published_at"`
}

func NewHeader() EventHeader {
	return EventHeader{
		ID:          uuid.NewString(),
		PublishedAt: time.Now().UTC(),
	}
}

type Event interface {
	isEvent()
}

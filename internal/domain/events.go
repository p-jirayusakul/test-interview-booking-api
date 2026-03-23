package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID
	Name          string
	MaxSeats      int
	WaitlistLimit int
	BookedCount   int
	WaitlistCount int
	Price         float64
	StartTime     time.Time
	EndTime       time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type EventsRepository interface {
	ListEvent(ctx context.Context) ([]Event, error)
}

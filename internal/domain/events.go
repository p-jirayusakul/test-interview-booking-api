package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	GetEvenById(ctx context.Context, id uuid.UUID) (Event, error)
	TXGetEventForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Event, error)
	TXExistsBooking(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error)
	TXUpdateEvent(ctx context.Context, tx *gorm.DB, payload Event) error
}

package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	orgerror "github.com/p-jirayusakul/test-interview-booking-api/pkg/errors"
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

type EventFilter struct {
	Name   string
	Page   int
	Limit  int
	Offset int
	Sort   string
	Order  string
}

type EventFilterResult struct {
	Items      []*Event
	Pagination EventFilterPagination
}

type EventFilterPagination struct {
	Page        int
	PageSize    int
	Total       int64
	HasNext     bool
	HasPrevious bool
}

type CreateEvent struct {
	Name          string
	MaxSeats      int
	WaitlistLimit int
	Price         float64
	StartTime     time.Time
	EndTime       time.Time
}

func (c *CreateEvent) Validate() error {

	if c.Name == "" {
		return orgerror.New(orgerror.CodeInvalidInput, "event name cannot be empty")
	}

	if c.MaxSeats < 1 {
		return orgerror.New(orgerror.CodeInvalidInput, "max seats cannot be less than 1")
	}

	if c.WaitlistLimit < 1 {
		return orgerror.New(orgerror.CodeInvalidInput, "waitlist limit cannot be less than 1")
	}

	if c.StartTime.After(c.EndTime) {
		return orgerror.New(orgerror.CodeInvalidInput, "start time cannot be after end time")
	}

	return nil
}

type EventsRepository interface {
	CreateEvent(ctx context.Context, payload *CreateEvent) error
	GetEvenById(ctx context.Context, id uuid.UUID) (*Event, error)
	TXGetEventForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*Event, error)
	TXExistsBooking(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error)
	TXUpdateEvent(ctx context.Context, tx *gorm.DB, payload *Event) error
	SearchEvents(ctx context.Context, payload *EventFilter) ([]*Event, error)
	CountEvents(ctx context.Context, search string) (int64, error)
}

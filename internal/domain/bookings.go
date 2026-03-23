package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusWaitlist  BookingStatus = "WAITLIST"
)

type Booking struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	UserID    uuid.UUID
	Status    BookingStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type BooksRepository interface {
	TXCreateBooking(ctx context.Context, tx *gorm.DB, payload Booking) error
}

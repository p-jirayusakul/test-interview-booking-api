package domain

import (
	"time"

	"github.com/google/uuid"
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
	UpdatedAt time.Time
}

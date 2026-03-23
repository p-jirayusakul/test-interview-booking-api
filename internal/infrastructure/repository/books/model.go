package books

import (
	"time"

	"github.com/google/uuid"
)

type bookingStatus string

const (
	bookingStatusConfirmed bookingStatus = "CONFIRMED"
	bookingStatusWaitlist  bookingStatus = "WAITLIST"
)

type booksRow struct {
	ID        uuid.UUID     `gorm:"column:id"`
	EventID   uuid.UUID     `gorm:"column:event_id"`
	UserID    uuid.UUID     `gorm:"column:user_id"`
	Status    bookingStatus `gorm:"column:status"`
	CreatedAt time.Time     `gorm:"column:created_at"`
	UpdatedAt *time.Time    `gorm:"column:updated_at"`
}

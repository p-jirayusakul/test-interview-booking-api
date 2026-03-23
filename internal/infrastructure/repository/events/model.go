package events

import (
	"time"

	"github.com/google/uuid"
)

type eventRow struct {
	ID            uuid.UUID  `gorm:"column:id"`
	Name          string     `gorm:"column:name"`
	MaxSeats      int        `gorm:"column:max_seats"`
	WaitlistLimit int        `gorm:"column:waitlist_limit"`
	BookedCount   int        `gorm:"column:booked_count"`
	WaitlistCount int        `gorm:"column:waitlist_count"`
	Price         float64    `gorm:"column:price"`
	StartTime     time.Time  `gorm:"column:start_time"`
	EndTime       time.Time  `gorm:"column:end_time"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

func whitelistSort(input string) string {
	switch input {
	case "name":
		return "name"
	case "maxSeats":
		return "max_seats"
	case "price":
		return "price"
	case "startTime":
		return "start_time"
	case "endTime":
		return "end_time"
	default:
		return "created_at"
	}
}

func whitelistOrder(input string) string {
	if input == "asc" {
		return "ASC"
	}
	return "DESC"
}

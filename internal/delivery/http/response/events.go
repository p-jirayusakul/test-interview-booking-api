package response

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID  `json:"id" example:"019d19e7-1626-7070-87aa-3879915fa117"`
	Name          string     `json:"name" example:"Go Workshop 101"`
	MaxSeats      int        `json:"maxSeats" example:"100"`
	WaitlistLimit int        `json:"waitlistLimit" example:"10"`
	BookedCount   int        `json:"bookedCount" example:"100"`
	WaitlistCount int        `json:"waitlistCount" example:"100"`
	Price         float64    `json:"price" example:"100000.00"`
	StartTime     time.Time  `json:"startTime" example:"2026-03-01T09:00:00+07:00"`
	EndTime       time.Time  `json:"endTime" example:"2026-03-01T12:00:00+07:00"`
	CreatedAt     time.Time  `json:"createdAt" example:"2026-03-01T09:00:00+07:00"`
	UpdatedAt     *time.Time `json:"updatedAt" example:"2026-03-01T12:00:00+07:00"`
}

type BookEvent struct {
	Status string `json:"status" example:"CONFIRMED"`
}

type SearchEventsResponse struct {
	Items      []Event                `json:"items"`
	Pagination SearchEventsPagination `json:"pagination"`
}

type SearchEventsPagination struct {
	Page        int   `json:"page" example:"1"`
	PageSize    int   `json:"pageSize" example:"10"`
	Total       int64 `json:"total" example:"10"`
	HasNext     bool  `json:"hasNext" example:"true"`
	HasPrevious bool  `json:"hasPrevious" example:"false"`
}

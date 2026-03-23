package response

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	MaxSeats      int        `json:"maxSeats"`
	WaitlistLimit int        `json:"waitlistLimit"`
	BookedCount   int        `json:"bookedCount"`
	WaitlistCount int        `json:"waitlistCount"`
	Price         float64    `json:"price"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       time.Time  `json:"endTime"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

type BookEvent struct {
	Status string `json:"status"`
}

type SearchEventsResponse struct {
	Items      []Event                `json:"items"`
	Pagination SearchEventsPagination `json:"pagination"`
}

type SearchEventsPagination struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	Total       int64 `json:"total"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

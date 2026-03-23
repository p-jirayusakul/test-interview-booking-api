package request

import (
	"time"
)

type SearchEventRequest struct {
	Search     string `query:"search" example:"Go Workshop 101"`
	PageNumber int    `query:"pageNumber" example:"1"`
	PageSize   int    `query:"pageSize" example:"10"`
	SortBy     string `query:"sortBy"`
	OrderBy    string `query:"orderBy"`
}

type CreateEvent struct {
	Name          string    `json:"name" example:"Go Workshop 101"`
	MaxSeats      int       `json:"maxSeats" example:"100"`
	WaitlistLimit int       `json:"waitlistLimit" example:"10"`
	Price         float64   `json:"price" example:"100000.00"`
	StartTime     time.Time `json:"startTime" example:"2026-03-01T09:00:00+07:00"`
	EndTime       time.Time `json:"endTime" example:"2026-03-01T12:00:00+07:00"`
}

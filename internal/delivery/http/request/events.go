package request

import (
	"time"
)

type SearchEventRequest struct {
	Search     string `query:"search"`
	PageNumber int    `query:"pageNumber"`
	PageSize   int    `query:"pageSize"`
	SortBy     string `query:"sortBy"`
	OrderBy    string `query:"orderBy"`
}

type CreateEvent struct {
	Name          string    `json:"name"`
	MaxSeats      int       `json:"maxSeats"`
	WaitlistLimit int       `json:"waitlistLimit"`
	Price         float64   `json:"price"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
}

package events

import "github.com/p-jirayusakul/test-interview-booking-api/internal/domain"

func mapEventRowsToDomain(rows []eventRow) []domain.Event {
	if rows == nil {
		return nil
	}

	events := make([]domain.Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, domain.Event{
			ID:            row.ID,
			Name:          row.Name,
			MaxSeats:      row.MaxSeats,
			WaitlistLimit: row.WaitlistLimit,
			BookedCount:   row.BookedCount,
			WaitlistCount: row.WaitlistCount,
			Price:         row.Price,
			StartTime:     row.StartTime,
			EndTime:       row.EndTime,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}

	return events
}

func mapEventRowToDomain(row eventRow) domain.Event {

	return domain.Event{
		ID:            row.ID,
		Name:          row.Name,
		MaxSeats:      row.MaxSeats,
		WaitlistLimit: row.WaitlistLimit,
		BookedCount:   row.BookedCount,
		WaitlistCount: row.WaitlistCount,
		Price:         row.Price,
		StartTime:     row.StartTime,
		EndTime:       row.EndTime,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

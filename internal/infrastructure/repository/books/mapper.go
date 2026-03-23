package books

import "github.com/p-jirayusakul/test-interview-booking-api/internal/domain"

func mapBookRowsToDomain(rows []booksRow) []domain.Booking {
	if rows == nil {
		return nil
	}

	books := make([]domain.Booking, 0, len(rows))
	for _, row := range rows {
		books = append(books, domain.Booking{
			ID:        row.ID,
			EventID:   row.EventID,
			UserID:    row.UserID,
			Status:    domain.BookingStatus(row.Status),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}

	return books
}

func mapBookRowToDomain(row booksRow) domain.Booking {
	return domain.Booking{
		ID:        row.ID,
		EventID:   row.EventID,
		UserID:    row.UserID,
		Status:    domain.BookingStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

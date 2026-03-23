package books

import (
	"context"

	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"gorm.io/gorm"
)

type booksRepo struct {
	db *gorm.DB
}

func NewBooksRepository(db *gorm.DB) domain.BooksRepository {
	return &booksRepo{db: db}
}

func (r *booksRepo) TXCreateBooking(ctx context.Context, tx *gorm.DB, payload domain.Booking) error {

	err := gorm.G[booksRow](tx).Exec(ctx, `
		insert into bookings (event_id, user_id, status) values ($1, $2, $3);
		`, payload.EventID, payload.UserID, payload.Status)
	if err != nil {
		return err
	}

	return nil
}

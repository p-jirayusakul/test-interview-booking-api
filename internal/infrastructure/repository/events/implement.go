package events

import (
	"context"

	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"gorm.io/gorm"
)

type eventsRepo struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) domain.EventsRepository {
	return &eventsRepo{db: db}
}

func (r *eventsRepo) ListEvent(ctx context.Context) ([]domain.Event, error) {
	result, err := gorm.G[eventRow](r.db).Raw(`
		SELECT
			id,
			name,
			max_seats,
			waitlist_limit,
			booked_count,
			waitlist_count,
			price,
			start_time,
			end_time,
			created_at,
			updated_at
		FROM public.events;
		`).Find(ctx)
	if err != nil {
		return nil, err
	}

	return mapEventRowsToDomain(result), nil
}

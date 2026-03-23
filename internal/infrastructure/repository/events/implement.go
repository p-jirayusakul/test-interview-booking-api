package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"gorm.io/gorm"
)

type eventsRepo struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) domain.EventsRepository {
	return &eventsRepo{db: db}
}

func (r *eventsRepo) CreateEvent(ctx context.Context, payload *domain.CreateEvent) error {

	err := gorm.G[eventRow](r.db).Exec(ctx, `
		INSERT INTO public.events (name, max_seats, waitlist_limit, price, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?);
		`, payload.Name, payload.MaxSeats, payload.WaitlistLimit, payload.Price, payload.StartTime, payload.EndTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventsRepo) SearchEvents(ctx context.Context, payload *domain.EventFilter) ([]*domain.Event, error) {
	sortCol := whitelistSort(payload.Sort)
	orderDir := whitelistOrder(payload.Order)

	var (
		conditions []string
		args       []interface{}
		argPos     = 1
	)

	if payload.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE '%%' || $%d || '%%'", argPos))
		args = append(args, payload.Name)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	args = append(args, payload.Limit, payload.Offset)

	query := fmt.Sprintf(`
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
		FROM public.events
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortCol, orderDir, argPos, argPos+1)

	result, err := gorm.G[eventRow](r.db).
		Raw(query, args...).
		Find(ctx)

	if err != nil {
		return nil, err
	}

	return mapEventRowsToDomain(result), nil
}

func (r *eventsRepo) CountEvents(ctx context.Context, search string) (int64, error) {

	result, err := gorm.G[int64](r.db).Raw(`
	SELECT COUNT(*)
			FROM public.events
			WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
	`, search).First(ctx)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *eventsRepo) GetEvenById(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
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
		FROM public.events WHERE id = ?;
		`, id).First(ctx)
	if err != nil {
		return nil, err
	}

	return mapEventRowToDomain(result), nil
}

func (r *eventsRepo) TXGetEventForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Event, error) {
	result, err := gorm.G[eventRow](tx).Raw(`
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
		FROM public.events WHERE id = ? FOR UPDATE;
		`, id).First(ctx)
	if err != nil {
		return nil, err
	}

	return mapEventRowToDomain(result), nil
}

func (r *eventsRepo) TXExistsBooking(ctx context.Context, tx *gorm.DB, eventId, userId uuid.UUID) (bool, error) {
	result, err := gorm.G[bool](tx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM public.bookings
			WHERE event_id = ? AND user_id = ?
		)
		`, eventId, userId).First(ctx)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (r *eventsRepo) TXUpdateEvent(ctx context.Context, tx *gorm.DB, payload *domain.Event) error {

	err := gorm.G[eventRow](tx).Exec(ctx, `
		UPDATE public.events
		SET
			booked_count = ?,
			waitlist_count = ?,
			updated_at = NOW()
		WHERE id = ?
		`, payload.BookedCount, payload.WaitlistCount, payload.ID)
	if err != nil {
		return err
	}

	return nil
}

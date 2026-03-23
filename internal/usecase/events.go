package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	orgerror "github.com/p-jirayusakul/test-interview-booking-api/pkg/errors"

	"gorm.io/gorm"
)

type EventsUseCase struct {
	txManager    domain.TransactionManager
	eventsRepo   domain.EventsRepository
	booksRepo    domain.BooksRepository
	timeLocation *time.Location
}

func NewEventsUseCase(eventsRepo domain.EventsRepository, booksRepo domain.BooksRepository, txManager domain.TransactionManager, timeLocation *time.Location) *EventsUseCase {
	return &EventsUseCase{eventsRepo: eventsRepo, booksRepo: booksRepo, txManager: txManager, timeLocation: timeLocation}
}

func (u *EventsUseCase) CreateEvent(ctx context.Context, payload *domain.CreateEvent) error {

	if err := payload.Validate(); err != nil {
		return err
	}

	err := u.eventsRepo.CreateEvent(ctx, payload)
	if err != nil {
		return orgerror.Wrap(orgerror.CodeSystem, "failed to create event", err)
	}

	return nil
}

func (u *EventsUseCase) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {

	if id == uuid.Nil {
		return nil, orgerror.New(orgerror.CodeInvalidInput, "event id cannot be nil")
	}

	result, err := u.eventsRepo.GetEvenById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, orgerror.New(orgerror.CodeEventNotFound, "event not found")
		}
		return nil, orgerror.Wrap(orgerror.CodeSystem, "failed to get event", err)
	}
	return result, nil
}

func (u *EventsUseCase) SearchEvents(ctx context.Context, payload *domain.EventFilter) (*domain.EventFilterResult, error) {

	if payload.Page < 1 {
		payload.Page = 1
	}

	if payload.Limit < 1 {
		payload.Limit = 10
	}

	if payload.Limit > 100 {
		payload.Limit = 100
	}

	payload.Offset = (payload.Page - 1) * payload.Limit

	var items []*domain.Event
	var total int64
	var itemsErr, totalErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		items, itemsErr = u.eventsRepo.SearchEvents(ctx, payload)
	}()

	go func() {
		defer wg.Done()
		total, totalErr = u.eventsRepo.CountEvents(ctx, payload.Name)
	}()

	wg.Wait()

	if itemsErr != nil {
		return nil, orgerror.Wrap(orgerror.CodeSystem, "failed to search event", itemsErr)
	}

	if totalErr != nil {
		return nil, orgerror.Wrap(orgerror.CodeSystem, "failed to count events", totalErr)
	}

	hasNext := payload.Page*payload.Limit < int(total)
	hasPrev := payload.Page > 1

	if total == 0 {
		hasNext = false
		hasPrev = false
	}

	return &domain.EventFilterResult{
		Items: items,
		Pagination: domain.EventFilterPagination{
			Page:        payload.Page,
			PageSize:    payload.Limit,
			Total:       total,
			HasNext:     hasNext,
			HasPrevious: hasPrev,
		},
	}, nil
}

func (u *EventsUseCase) BookEvent(ctx context.Context, eventId, userId uuid.UUID) (string, error) {

	var resultStatus string
	err := u.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {

		// 1. lock event by id
		event, err := u.eventsRepo.TXGetEventForUpdate(ctx, tx, eventId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return orgerror.New(orgerror.CodeEventNotFound, "event not found")
			}
			return orgerror.Wrap(orgerror.CodeSystem, "failed to get event for update", err)
		}

		// 2. check booking window
		now := time.Now().In(u.timeLocation)
		eventStartTime := event.StartTime.In(u.timeLocation)
		eventEndTime := event.EndTime.In(u.timeLocation)
		if now.Before(eventStartTime) || now.After(eventEndTime) {
			return orgerror.New(orgerror.CodeBookingClosed, "booking window closed")
		}

		// 3. check duplicate booking
		exists, err := u.eventsRepo.TXExistsBooking(ctx, tx, eventId, userId)
		if err != nil {
			return err
		}
		if exists {
			return orgerror.New(orgerror.CodeAlreadyBooked, "you have already booked this event")
		}

		// 4. decide status
		var status domain.BookingStatus
		if event.BookedCount < event.MaxSeats {
			status = domain.BookingStatusConfirmed
			event.BookedCount++
		} else if event.WaitlistCount < event.WaitlistLimit {
			status = domain.BookingStatusWaitlist
			event.WaitlistCount++
		} else {
			return orgerror.New(orgerror.CodeEventFull, "event is full")
		}

		// 5. create booking
		err = u.booksRepo.TXCreateBooking(ctx, tx, domain.Booking{
			EventID: eventId,
			UserID:  userId,
			Status:  status,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return orgerror.New(orgerror.CodeAlreadyBooked, "you have already booked this event")
			}
			return orgerror.Wrap(orgerror.CodeSystem, "failed to create booking", err)
		}

		// 6. update event counter
		err = u.eventsRepo.TXUpdateEvent(ctx, tx, event)
		if err != nil {
			return orgerror.Wrap(orgerror.CodeSystem, "failed to update event", err)
		}

		// 7. set result
		resultStatus = string(status)

		return nil
	})
	if err != nil {
		return "", err
	}

	return resultStatus, nil
}

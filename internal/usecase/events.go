package usecase

import (
	"context"

	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
)

type EventsUseCase struct {
	eventsRepo domain.EventsRepository
}

func NewEventsUseCase(eventsRepo domain.EventsRepository) *EventsUseCase {
	return &EventsUseCase{eventsRepo: eventsRepo}
}

func (e *EventsUseCase) ListEvent(ctx context.Context) ([]domain.Event, error) {
	return e.eventsRepo.ListEvent(ctx)
}

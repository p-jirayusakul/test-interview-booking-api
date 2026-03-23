package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http/response"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/usecase"
)

type EventsHandler struct {
	useCase *usecase.EventsUseCase
}

func NewEventsHandler(useCase *usecase.EventsUseCase) *EventsHandler {
	return &EventsHandler{useCase: useCase}
}

func (h *EventsHandler) ListEvent(c *echo.Context) error {
	result, err := h.useCase.ListEvent(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mapEventsToResponse(result))
}

func mapEventsToResponse(events []domain.Event) []response.Event {
	result := make([]response.Event, len(events))
	for i, event := range events {
		result[i] = response.Event{
			ID:            event.ID,
			Name:          event.Name,
			MaxSeats:      event.MaxSeats,
			WaitlistLimit: event.WaitlistLimit,
			BookedCount:   event.BookedCount,
			WaitlistCount: event.WaitlistCount,
			Price:         event.Price,
			StartTime:     event.StartTime,
			EndTime:       event.EndTime,
			CreatedAt:     event.CreatedAt,
			UpdatedAt:     event.UpdatedAt,
		}
	}
	return result
}

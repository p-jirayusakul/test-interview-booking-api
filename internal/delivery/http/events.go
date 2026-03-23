package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http/request"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http/response"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/usecase"
	orgerror "github.com/p-jirayusakul/test-interview-booking-api/pkg/errors"
	orgresponse "github.com/p-jirayusakul/test-interview-booking-api/pkg/response"
)

type EventsHandler struct {
	useCase *usecase.EventsUseCase
}

func NewEventsHandler(useCase *usecase.EventsUseCase) *EventsHandler {
	return &EventsHandler{useCase: useCase}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param event body request.CreateEvent true "Event details"
// @Success 201 {object} orgresponse.Response[any]
// @Failure 400 {object} orgresponse.Response[any]
// @Failure 500 {object} orgresponse.Response[any]
// @Router /api/v1/events [post]
func (h *EventsHandler) CreateEvent(c *echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)

	var createEvent request.CreateEvent
	err := c.Bind(&createEvent)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	payload := domain.CreateEvent{
		Name:          createEvent.Name,
		MaxSeats:      createEvent.MaxSeats,
		WaitlistLimit: createEvent.WaitlistLimit,
		Price:         createEvent.Price,
		StartTime:     createEvent.StartTime,
		EndTime:       createEvent.EndTime,
	}

	err = h.useCase.CreateEvent(c.Request().Context(), payload)
	if err != nil {
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	return c.JSON(http.StatusCreated, orgresponse.Response[any]{
		RequestID: requestId,
		Success:   true,
	})
}

// GetEvent godoc
// @Summary Get event by ID
// @Description Retrieve event details by event ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} orgresponse.Response[response.Event]
// @Failure 400 {object} orgresponse.Response[any]
// @Failure 404 {object} orgresponse.Response[any]
// @Failure 500 {object} orgresponse.Response[any]
// @Router /api/v1/events/{id} [get]
func (h *EventsHandler) GetEvent(c *echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)

	eventIdStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	result, err := h.useCase.GetEvent(c.Request().Context(), eventId)
	if err != nil {
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	return c.JSON(http.StatusOK, orgresponse.Response[response.Event]{
		RequestID: requestId,
		Success:   true,
		Data:      new(mapEventResponse(result)),
	})
}

// SearchEvents godoc
// @Summary Search events
// @Description Search and filter events with pagination and sorting
// @Tags events
// @Accept json
// @Produce json
// @Param search query string false "Search by event name"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param sortBy query string false "Sort by field"
// @Param orderBy query string false "Order direction (asc/desc)"
// @Success 200 {object} orgresponse.Response[response.SearchEventsResponse]
// @Failure 400 {object} orgresponse.Response[any]
// @Failure 500 {object} orgresponse.Response[any]
// @Router /api/v1/events/search [get]
func (h *EventsHandler) SearchEvents(c *echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)

	var searchEventRequest request.SearchEventRequest
	err := c.Bind(&searchEventRequest)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	payload := domain.EventFilter{
		Name:  searchEventRequest.Search,
		Page:  searchEventRequest.PageNumber,
		Limit: searchEventRequest.PageSize,
		Sort:  searchEventRequest.SortBy,
		Order: searchEventRequest.OrderBy,
	}

	result, err := h.useCase.SearchEvents(c.Request().Context(), payload)
	if err != nil {
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	return c.JSON(http.StatusOK, orgresponse.Response[response.SearchEventsResponse]{
		RequestID: requestId,
		Success:   true,
		Data:      new(mapSearchEventResponse(result)),
	})
}

// BookEvent godoc
// @Summary Book an event
// @Description Book a seat for an event or join the waitlist
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param X-User-Id header string true "User ID"
// @Success 201 {object} orgresponse.Response[response.BookEvent]
// @Failure 400 {object} orgresponse.Response[any]
// @Failure 401 {object} orgresponse.Response[any]
// @Failure 404 {object} orgresponse.Response[any]
// @Failure 500 {object} orgresponse.Response[any]
// @Router /api/v1/events/{id}/book [post]
func (h *EventsHandler) BookEvent(c *echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	userIdStr := c.Request().Header.Get("X-User-Id")
	if userIdStr == "" {
		err := orgerror.New(orgerror.CodeUnauthorized, "user id is required")
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	eventIdStr, err := echo.PathParam[string](c, "id")
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	result, err := h.useCase.BookEvent(c.Request().Context(), eventId, userId)
	if err != nil {
		return c.JSON(orgerror.HTTPStatus(err), orgresponse.ErrorResponse(err, requestId))
	}

	return c.JSON(http.StatusCreated, orgresponse.Response[response.BookEvent]{
		RequestID: requestId,
		Success:   true,
		Data:      &response.BookEvent{Status: result},
	})
}

func mapSearchEventResponse(payload domain.EventFilterResult) response.SearchEventsResponse {
	return response.SearchEventsResponse{
		Items: mapSearchEventResponseItems(payload.Items),
		Pagination: response.SearchEventsPagination{
			Total:       payload.Pagination.Total,
			Page:        payload.Pagination.Page,
			PageSize:    payload.Pagination.PageSize,
			HasNext:     payload.Pagination.HasNext,
			HasPrevious: payload.Pagination.HasPrevious,
		},
	}
}

func mapSearchEventResponseItems(rows []domain.Event) []response.Event {
	if rows == nil {
		return nil
	}

	events := make([]response.Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, response.Event{
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

func mapEventResponse(row domain.Event) response.Event {
	return response.Event{
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

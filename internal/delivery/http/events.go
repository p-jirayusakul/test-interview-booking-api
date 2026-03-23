package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http/response"
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

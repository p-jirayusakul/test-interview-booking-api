package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	orgerror "github.com/p-jirayusakul/test-interview-booking-api/pkg/errors"
)

func BindEventsRoutes(routesGroup *echo.Group, handler *EventsHandler) {
	eventsRouter := routesGroup.Group("/events")
	eventsRouter.POST("", handler.CreateEvent)
	eventsRouter.GET("/:id", handler.GetEvent)
	eventsRouter.GET("/search", handler.SearchEvents)
	eventsRouter.POST("/:id/book", handler.BookEvent)
}

func convertStringToUUID(c *echo.Context, s string) (uuid.UUID, error) {

	str, err := echo.PathParam[string](c, s)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return uuid.Nil, err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		err = orgerror.New(orgerror.CodeInvalidInput, err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

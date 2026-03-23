package http

import (
	"github.com/labstack/echo/v5"
)

func BindEventsRoutes(routesGroup *echo.Group, handler *EventsHandler) {
	eventsRouter := routesGroup.Group("/events")
	eventsRouter.GET("", handler.ListEvent)
}

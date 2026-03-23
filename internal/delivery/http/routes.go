package http

import (
	"github.com/labstack/echo/v5"
)

func BindEventsRoutes(routesGroup *echo.Group, handler *EventsHandler) {
	eventsRouter := routesGroup.Group("/events")
	eventsRouter.POST("", handler.CreateEvent)
	eventsRouter.GET("/:id", handler.GetEvent)
	eventsRouter.GET("/search", handler.SearchEvents)
	eventsRouter.POST("/:id/book", handler.BookEvent)
}

package bootstrap

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	orghttp "github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/books"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/events"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/postgres"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/usecase"
	"github.com/p-jirayusakul/test-interview-booking-api/pkg/logs"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newEchoServer(logger *zap.Logger) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(requestLoggerMiddleware(logger))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health/liveness", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}

func requestLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		Skipper: func(c *echo.Context) bool {
			return logs.ShouldSkip(c.Request().URL.Path)
		},
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			requestId := c.Response().Header().Get(echo.HeaderXRequestID)
			fields := []zap.Field{
				zap.String("request_id", requestId),
				zap.String("method", c.Request().Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
			}
			userIdStr := c.Request().Header.Get("X-User-Id")
			if userIdStr != "" {
				fields = append(fields, zap.String("user_id", userIdStr))
			}

			logger.Info("HTTP INBOUND REQUEST", fields...)
			return nil
		},
	})
}

func initEventsHandler(routesGroup *echo.Group, dbConn *gorm.DB, timeLocation *time.Location) {
	eventsRepo := events.NewEventsRepository(dbConn)
	booksRepo := books.NewBooksRepository(dbConn)
	txManager := postgres.NewTransactionManager(dbConn)
	eventsUseCase := usecase.NewEventsUseCase(eventsRepo, booksRepo, txManager, timeLocation)
	eventsHandler := orghttp.NewEventsHandler(eventsUseCase)
	orghttp.BindEventsRoutes(routesGroup, eventsHandler)
}

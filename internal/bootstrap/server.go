package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	orghttp "github.com/p-jirayusakul/test-interview-booking-api/internal/delivery/http"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/config"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/books"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/events"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/postgres"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/usecase"
	"github.com/p-jirayusakul/test-interview-booking-api/pkg/logs"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
	"go.uber.org/zap"

	_ "github.com/p-jirayusakul/test-interview-booking-api/docs"
	"gorm.io/gorm"
)

func NewServer() (*http.Server, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	dbConn, err := postgres.NewConnection(postgres.Config{
		Host:     cfg.DatabaseCfg.Host,
		Port:     cfg.DatabaseCfg.Port,
		User:     cfg.DatabaseCfg.User,
		Password: cfg.DatabaseCfg.Password,
		DBName:   cfg.DatabaseCfg.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	routes := echo.New()
	routes.Use(middleware.Recover())
	routes.Use(middleware.RequestID())
	routes.GET("/swagger/*", echoSwagger.WrapHandler)
	routes.GET("/health/liveness", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	logger, err := logs.InitLog(cfg.AppCfg.Env)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	routes.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	}))

	routesGroup := routes.Group(cfg.AppCfg.BaseURL)

	loc, err := time.LoadLocation(cfg.AppCfg.TZ)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	initEventsHandler(routesGroup, dbConn, loc)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppCfg.Port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("Server started on port: ", cfg.AppCfg.Port)
	return server, nil
}

func initEventsHandler(routesGroup *echo.Group, dbConn *gorm.DB, timeLocation *time.Location) {
	eventsRepo := events.NewEventsRepository(dbConn)
	booksRepo := books.NewBooksRepository(dbConn)
	txManager := postgres.NewTransactionManager(dbConn)
	eventsUseCase := usecase.NewEventsUseCase(eventsRepo, booksRepo, txManager, timeLocation)
	eventsHandler := orghttp.NewEventsHandler(eventsUseCase)
	orghttp.BindEventsRoutes(routesGroup, eventsHandler)
}

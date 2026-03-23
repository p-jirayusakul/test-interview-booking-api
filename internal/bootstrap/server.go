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

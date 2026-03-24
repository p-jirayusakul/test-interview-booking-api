package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/p-jirayusakul/test-interview-booking-api/docs"
	"github.com/p-jirayusakul/test-interview-booking-api/pkg/logs"
	"gorm.io/gorm"
)

type Server struct {
	HttpServer *http.Server
	DB         *gorm.DB
}

func NewServer() (*Server, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	dbConn, err := initDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	logger, err := logs.InitLog(cfg.AppCfg.Env)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	routes := newEchoServer(logger)
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
	return &Server{
		DB:         dbConn,
		HttpServer: server,
	}, nil
}

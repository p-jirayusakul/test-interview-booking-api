package bootstrap

import (
	"fmt"

	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/config"
	"github.com/p-jirayusakul/test-interview-booking-api/internal/infrastructure/repository/postgres"
	"gorm.io/gorm"
)

func initConfig() (*config.Config, error) {
	return config.InitConfig()
}

func initDatabase(cfg *config.Config) (dbConn *gorm.DB, err error) {
	dbConn, err = postgres.NewConnection(postgres.Config{
		Host:     cfg.DatabaseCfg.Host,
		Port:     cfg.DatabaseCfg.Port,
		User:     cfg.DatabaseCfg.User,
		Password: cfg.DatabaseCfg.Password,
		DBName:   cfg.DatabaseCfg.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	return dbConn, nil
}

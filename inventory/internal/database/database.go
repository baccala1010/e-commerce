package database

import (
	"time"

	"github.com/baccala1010/e-commerce/inventory/internal/config"
	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/pkg/postgre"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection and performs migrations
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.Logging.Level == "debug" {
		logLevel = logger.Info
	}

	// Create DB options
	dbOptions := postgre.NewDBOptions().
		WithHost(cfg.Database.Host).
		WithPort(cfg.Database.Port).
		WithDatabase(cfg.Database.Name).
		WithUsername(cfg.Database.Username).
		WithPassword(cfg.Database.Password).
		WithSSLMode(cfg.Database.SSLMode).
		WithMaxIdleConnections(cfg.Database.MaxIdleConnections).
		WithMaxOpenConnections(cfg.Database.MaxOpenConnections).
		WithConnectionMaxLifetime(parseDuration(cfg.Database.ConnectionMaxLifetime))

	// Connect to database
	db, err := postgre.Connect(dbOptions, logLevel)
	if err != nil {
		return nil, err
	}

	// Auto-migrate the models
	if err := db.AutoMigrate(
		&model.Product{},
		&model.Category{},
		&model.Discount{},
	); err != nil {
		return nil, err
	}

	return db.GetConnection(), nil
}

// parseDuration parses a duration string with a fallback to 1 hour
func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logrus.Warnf("Invalid duration string: %s, using default of 1 hour", durationStr)
		return time.Hour
	}
	return duration
}

package database

import (
	"fmt"

	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/baccala1010/e-commerce/statistics/pkg/postgre"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB creates and initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.Database.Debug {
		logLevel = logger.Info
	}

	dbOptions := postgre.NewDBOptions().
		WithHost(cfg.Database.Host).
		WithPort(cfg.Database.Port).
		WithDatabase(cfg.Database.DBName).
		WithUsername(cfg.Database.User).
		WithPassword(cfg.Database.Password).
		WithSSLMode(cfg.Database.SSLMode)

	db, err := postgre.Connect(dbOptions, logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&model.UserOrderStatistic{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed initial data if needed
	if err := seedInitialData(db); err != nil {
		return nil, fmt.Errorf("failed to seed initial data: %w", err)
	}

	return db, nil
}

// seedInitialData adds some initial test data if the tables are empty
func seedInitialData(db *gorm.DB) error {
	// Check if we already have data
	var count int64
	if err := db.Model(&model.UserOrderStatistic{}).Count(&count).Error; err != nil {
		return err
	}

	// If we have data, skip seeding
	if count > 0 {
		return nil
	}

	// Add a test record for development and testing
	testUserID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // This is the example UUID used in test files
	testStat := model.UserOrderStatistic{
		UserID:     testUserID,
		OrderCount: 5,
	}

	return db.Create(&testStat).Error
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

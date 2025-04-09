package postgre

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB represents a PostgreSQL database connection
type DB struct {
	conn *gorm.DB
}

// Connect establishes a connection to the PostgreSQL database
func Connect(options *DBOptions, logLevel logger.LogLevel) (*DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	conn, err := gorm.Open(postgres.Open(options.ConnectionString()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool parameters
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(options.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(options.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(options.ConnectionMaxLifetime)

	logrus.Info("Successfully connected to PostgreSQL database")
	return &DB{conn: conn}, nil
}

// GetConnection returns the underlying GORM database connection
func (db *DB) GetConnection() *gorm.DB {
	return db.conn
}

// AutoMigrate automatically migrates the schemas for the given models
func (db *DB) AutoMigrate(models ...interface{}) error {
	if err := db.conn.AutoMigrate(models...); err != nil {
		logrus.Errorf("Failed to migrate database: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logrus.Info("Successfully migrated database schema")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	logrus.Info("Successfully closed database connection")
	return nil
}

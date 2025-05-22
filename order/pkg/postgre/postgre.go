package postgre

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB represents a PostgreSQL database connection
type DB struct {
	conn *gorm.DB
}

// EnsureDatabaseExists checks if the database exists and creates it if it doesn't
func EnsureDatabaseExists(options *DBOptions) error {
	// Connect to the default "postgres" database to create our target database
	defaultOptions := NewDBOptions().
		WithHost(options.Host).
		WithPort(options.Port).
		WithDatabase("postgres"). // Connect to default postgres database
		WithUsername(options.Username).
		WithPassword(options.Password).
		WithSSLMode(options.SSLMode)

	// Connect to the default database
	connectionString := defaultOptions.ConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to default database: %w", err)
	}
	defer db.Close()

	// Check if our target database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, options.Database).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create the database if it doesn't exist
	if !exists {
		logrus.Infof("Database %s does not exist, creating it now", options.Database)
		_, err = db.Exec("CREATE DATABASE " + options.Database)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		logrus.Infof("Successfully created database %s", options.Database)
	}

	return nil
}

// Connect establishes a connection to the PostgreSQL database with retries
func Connect(options *DBOptions, logLevel logger.LogLevel) (*DB, error) {
	// First, make sure the database exists
	if err := EnsureDatabaseExists(options); err != nil {
		logrus.Warnf("Failed to ensure database exists: %v. Will try to connect anyway.", err)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Retry logic for database connection
	var conn *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logrus.Infof("Attempting to connect to database (attempt %d/%d)", attempt, maxRetries)
		conn, err = gorm.Open(postgres.Open(options.ConnectionString()), gormConfig)
		if err == nil {
			break // Successfully connected
		}

		logrus.Warnf("Failed to connect to database (attempt %d/%d): %v", attempt, maxRetries, err)
		if attempt < maxRetries {
			logrus.Infof("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
			// Increase retry interval for next attempt
			retryInterval *= 2
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
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

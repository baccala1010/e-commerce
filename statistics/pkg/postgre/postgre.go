package postgre

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(options *DBOptions, logLevel logger.LogLevel) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}
	conn, err := gorm.Open(postgres.Open(options.ConnectionString()), gormConfig)
	if err != nil {
		return nil, err
	}
	// Set connection pool parameters
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(options.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(options.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(options.ConnectionMaxLifetime)
	return conn, nil
}

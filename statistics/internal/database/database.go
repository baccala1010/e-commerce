package database

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/baccala1010/e-commerce/statistics/pkg/postgre"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, dsn string) (*Database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	// You can add a debug flag to config if needed
	// if cfg.Logging.Level == "debug" { logLevel = logger.Info }

	dbOptions := postgre.NewDBOptions().
		WithHost(cfg.Database.Host).
		WithPort(cfg.Database.Port).
		WithDatabase(cfg.Database.DBName).
		WithUsername(cfg.Database.User).
		WithPassword(cfg.Database.Password).
		WithSSLMode("disable") // or cfg.Database.SSLMode if present

	db, err := postgre.Connect(dbOptions, logLevel)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.UserOrderStatistic{}); err != nil {
		return nil, err
	}

	return db, nil
}

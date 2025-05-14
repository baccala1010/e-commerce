package postgre

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Options represents PostgreSQL connection options
type Options struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	SSLMode        string
	MaxConnections int
}

// New creates a new PostgreSQL connection pool
func New(ctx context.Context, opts Options) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d",
		opts.Host,
		opts.Port,
		opts.User,
		opts.Password,
		opts.Database,
		opts.SSLMode,
		opts.MaxConnections,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return pool, nil
}

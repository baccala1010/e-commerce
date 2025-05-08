package database

import (
	"context"
	"fmt"
	"log"

	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/pkg/postgre"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDatabase initializes the PostgreSQL database connection
func InitDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	opts := postgre.Options{
		Host:           cfg.DB.Host,
		Port:           cfg.DB.Port,
		User:           cfg.DB.User,
		Password:       cfg.DB.Password,
		Database:       cfg.DB.Name,
		SSLMode:        cfg.DB.SSLMode,
		MaxConnections: cfg.DB.MaxConnections,
	}

	pool, err := postgre.New(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	// Create necessary tables if they don't exist
	if err := createTables(ctx, pool); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return pool, nil
}

// createTables creates all necessary tables if they don't exist
func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	// Users statistics table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT,
			name TEXT,
			registration_date TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Orders statistics table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			total_amount DECIMAL(10, 2) NOT NULL,
			order_status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	// Products statistics table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			category_id TEXT,
			price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Order items table to track products in orders
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id TEXT NOT NULL,
			product_id TEXT NOT NULL,
			quantity INT NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create order_items table: %w", err)
	}

	log.Println("Database tables created successfully")
	return nil
}
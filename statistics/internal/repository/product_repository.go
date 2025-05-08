package repository

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository creates a new PostgreSQL implementation of ProductRepository
func NewProductRepository(pool *pgxpool.Pool) ProductRepository {
	return &productRepository{
		pool: pool,
	}
}

// Create inserts a new product record
func (r *productRepository) Create(ctx context.Context, product model.Product) error {
	query := `
		INSERT INTO products (id, name, category_id, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			category_id = $3,
			price = $4,
			updated_at = $6
	`

	_, err := r.pool.Exec(ctx, query,
		product.ID,
		product.Name,
		product.CategoryID,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// Update updates an existing product record
func (r *productRepository) Update(ctx context.Context, product model.Product) error {
	query := `
		UPDATE products
		SET name = $2, category_id = $3, price = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		product.ID,
		product.Name,
		product.CategoryID,
		product.Price,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// FindByID retrieves a product by ID
func (r *productRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	query := `
		SELECT id, name, category_id, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product model.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.CategoryID,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &product, nil
}
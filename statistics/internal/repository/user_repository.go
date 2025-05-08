package repository

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL implementation of UserRepository
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		pool: pool,
	}
}

// Create inserts a new user record
func (r *userRepository) Create(ctx context.Context, user model.User) error {
	query := `
		INSERT INTO users (id, email, name, registration_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.RegistrationDate,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user record
func (r *userRepository) Update(ctx context.Context, user model.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, name, registration_date, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.RegistrationDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// CountAll counts all users
func (r *userRepository) CountAll(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM users"

	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
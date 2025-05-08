package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new PostgreSQL implementation of OrderRepository
func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return &orderRepository{
		pool: pool,
	}
}

// Create inserts a new order record
func (r *orderRepository) Create(ctx context.Context, order model.Order) error {
	query := `
		INSERT INTO orders (id, user_id, total_amount, order_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			total_amount = $3,
			order_status = $4,
			updated_at = $6
	`

	_, err := r.pool.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.TotalAmount,
		order.OrderStatus,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// Update updates an existing order record
func (r *orderRepository) Update(ctx context.Context, order model.Order) error {
	query := `
		UPDATE orders
		SET total_amount = $2, order_status = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		order.ID,
		order.TotalAmount,
		order.OrderStatus,
		order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// AddOrderItems adds items to an order
func (r *orderRepository) AddOrderItems(ctx context.Context, orderID string, items []model.OrderItem) error {
	for _, item := range items {
		query := `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING
		`

		_, err := r.pool.Exec(ctx, query,
			orderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)

		if err != nil {
			return fmt.Errorf("failed to add order item: %w", err)
		}
	}

	return nil
}

// FindByID retrieves an order by ID
func (r *orderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	query := `
		SELECT id, user_id, total_amount, order_status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order model.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.OrderStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Get order items
	itemsQuery := `
		SELECT id, product_id, quantity, price, created_at
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.pool.Query(ctx, itemsQuery, id)
	if err != nil {
		log.Printf("failed to get order items: %v", err)
	} else {
		defer rows.Close()

		var items []model.OrderItem
		for rows.Next() {
			var item model.OrderItem
			err := rows.Scan(
				&item.ID,
				&item.ProductID,
				&item.Quantity,
				&item.Price,
				&item.CreatedAt,
			)
			if err != nil {
				log.Printf("failed to scan order item: %v", err)
				continue
			}
			item.OrderID = id
			items = append(items, item)
		}
		order.Items = items
	}

	return &order, nil
}

// FindByUserID retrieves all orders for a user
func (r *orderRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	query := `
		SELECT id, user_id, total_amount, order_status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by user ID: %w", err)
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.OrderStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetUserOrdersStatistics calculates statistics for a user's orders
func (r *orderRepository) GetUserOrdersStatistics(ctx context.Context, userID string) (*model.UserOrderStatistics, error) {
	query := `
		SELECT 
			COUNT(*) as total_orders,
			SUM(total_amount) as total_spent,
			MIN(created_at) as first_order_at,
			MAX(created_at) as last_order_at
		FROM orders
		WHERE user_id = $1
	`

	var stats model.UserOrderStatistics
	var totalOrders int
	var totalSpent float64
	var firstOrderAt, lastOrderAt *time.Time // Using nullable time to handle case with no orders

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&totalOrders,
		&totalSpent,
		&firstOrderAt,
		&lastOrderAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get order statistics: %w", err)
	}

	stats.UserID = userID
	stats.TotalOrders = totalOrders
	stats.TotalSpent = totalSpent

	if totalOrders > 0 {
		stats.AverageOrderValue = totalSpent / float64(totalOrders)
		if firstOrderAt != nil {
			stats.FirstOrderAt = *firstOrderAt
		}
		if lastOrderAt != nil {
			stats.LastOrderAt = *lastOrderAt
		}

		// Get hourly distribution
		timeDistribution, err := r.GetHourlyDistribution(ctx, userID)
		if err != nil {
			log.Printf("failed to get order time distribution: %v", err)
		} else {
			stats.OrderTimeDistribution = timeDistribution
		}
	}

	return &stats, nil
}

// GetHourlyDistribution gets the hourly distribution of orders for a user
func (r *orderRepository) GetHourlyDistribution(ctx context.Context, userID string) ([]model.OrderTimeOfDay, error) {
	query := `
		SELECT 
			to_char(created_at, 'HH24:00-HH24:59') as hour,
			COUNT(*) as order_count
		FROM orders
		WHERE user_id = $1
		GROUP BY hour
		ORDER BY hour
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hourly distribution: %w", err)
	}
	defer rows.Close()

	var distribution []model.OrderTimeOfDay
	for rows.Next() {
		var item model.OrderTimeOfDay
		err := rows.Scan(&item.Hour, &item.OrderCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hourly distribution: %w", err)
		}
		distribution = append(distribution, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating hourly distribution: %w", err)
	}

	return distribution, nil
}

// Time is a wrapper around time.Time for handling NULL times from the database
type Time struct {
	model.Time
	Valid bool
}

// Scan implements the Scanner interface for Time
func (t *Time) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	return nil
}
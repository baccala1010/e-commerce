package repository

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
)

// UserRepository defines methods for user statistics data access
type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	CountAll(ctx context.Context) (int, error)
}

// OrderRepository defines methods for order statistics data access
type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Update(ctx context.Context, order model.Order) error
	AddOrderItems(ctx context.Context, orderID string, items []model.OrderItem) error
	FindByID(ctx context.Context, id string) (*model.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	GetUserOrdersStatistics(ctx context.Context, userID string) (*model.UserOrderStatistics, error)
	GetHourlyDistribution(ctx context.Context, userID string) ([]model.OrderTimeOfDay, error)
}

// ProductRepository defines methods for product statistics data access
type ProductRepository interface {
	Create(ctx context.Context, product model.Product) error
	Update(ctx context.Context, product model.Product) error
	FindByID(ctx context.Context, id string) (*model.Product, error)
}
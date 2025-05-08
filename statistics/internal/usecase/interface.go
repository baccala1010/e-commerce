package usecase

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
)

// StatisticsUsecase defines the use cases for statistics functionality
type StatisticsUsecase interface {
	// GetUserOrdersStatistics gets statistics about a user's orders
	GetUserOrdersStatistics(ctx context.Context, userID string) (*model.UserOrderStatistics, error)
	
	// GetUserStatistics gets general user statistics
	GetUserStatistics(ctx context.Context) (*model.UserStatistics, error)
}
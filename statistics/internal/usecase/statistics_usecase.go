package usecase

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/baccala1010/e-commerce/statistics/internal/repository"
)

type statisticsUsecase struct {
	userRepo  repository.UserRepository
	orderRepo repository.OrderRepository
}

// NewStatisticsUsecase creates a new statistics use case implementation
func NewStatisticsUsecase(userRepo repository.UserRepository, orderRepo repository.OrderRepository) StatisticsUsecase {
	return &statisticsUsecase{
		userRepo:  userRepo,
		orderRepo: orderRepo,
	}
}

// GetUserOrdersStatistics retrieves statistics about a user's orders
func (u *statisticsUsecase) GetUserOrdersStatistics(ctx context.Context, userID string) (*model.UserOrderStatistics, error) {
	// Check if user exists
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get order statistics for the user
	stats, err := u.orderRepo.GetUserOrdersStatistics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order statistics: %w", err)
	}

	return stats, nil
}

// GetUserStatistics retrieves general user statistics
func (u *statisticsUsecase) GetUserStatistics(ctx context.Context) (*model.UserStatistics, error) {
	// Count total registered users
	totalUsers, err := u.userRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	stats := &model.UserStatistics{
		TotalRegisteredUsers: totalUsers,
	}

	return stats, nil
}

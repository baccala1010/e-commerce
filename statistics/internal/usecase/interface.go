package usecase

import (
	"context"

	"github.com/google/uuid"
)

// UserStatistic represents statistics for a single user
type UserStatistic struct {
	UserID     string
	OrderCount int
}

type StatisticsUsecase interface {
	GetUserOrderCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetAllUserOrderStatistics(ctx context.Context, page, pageSize int) ([]UserStatistic, int64, error)
}

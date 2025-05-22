package usecase

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/repository"
	"github.com/google/uuid"
)

type statisticsUsecase struct {
	repo repository.StatisticsRepository
}

func NewStatisticsUsecase(repo repository.StatisticsRepository) StatisticsUsecase {
	return &statisticsUsecase{repo: repo}
}

func (u *statisticsUsecase) GetUserOrderCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return u.repo.GetOrderCount(ctx, userID)
}

func (u *statisticsUsecase) GetAllUserOrderStatistics(ctx context.Context, page, pageSize int) ([]UserStatistic, int64, error) {
	stats, total, err := u.repo.GetAllStatistics(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Map from model to usecase type
	result := make([]UserStatistic, len(stats))
	for i, stat := range stats {
		result[i] = UserStatistic{
			UserID:     stat.UserID.String(),
			OrderCount: stat.OrderCount,
		}
	}

	return result, total, nil
}

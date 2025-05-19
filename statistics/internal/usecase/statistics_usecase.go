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

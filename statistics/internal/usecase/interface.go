package usecase

import (
	"context"

	"github.com/google/uuid"
)

type StatisticsUsecase interface {
	GetUserOrderCount(ctx context.Context, userID uuid.UUID) (int, error)
}

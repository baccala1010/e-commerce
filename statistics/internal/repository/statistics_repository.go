package repository

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatisticsRepository interface {
	IncrementOrderCount(ctx context.Context, userID uuid.UUID) error
	DecrementOrderCount(ctx context.Context, userID uuid.UUID) error
	GetOrderCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetAllStatistics(ctx context.Context, page, pageSize int) ([]model.UserOrderStatistic, int64, error)
}

type statisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) IncrementOrderCount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.UserOrderStatistic{}).
		Where("user_id = ?", userID).
		Assign(map[string]interface{}{"order_count": gorm.Expr("order_count + 1")}).
		FirstOrCreate(&model.UserOrderStatistic{UserID: userID}).Error
}

func (r *statisticsRepository) DecrementOrderCount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.UserOrderStatistic{}).
		Where("user_id = ?", userID).
		UpdateColumn("order_count", gorm.Expr("GREATEST(order_count - 1, 0)")).Error
}

func (r *statisticsRepository) GetOrderCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var stat model.UserOrderStatistic
	err := r.db.WithContext(ctx).First(&stat, "user_id = ?", userID).Error
	if err != nil {
		return 0, err
	}
	return stat.OrderCount, nil
}

func (r *statisticsRepository) GetAllStatistics(ctx context.Context, page, pageSize int) ([]model.UserOrderStatistic, int64, error) {
	var stats []model.UserOrderStatistic
	var totalCount int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&model.UserOrderStatistic{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Query with pagination
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	query := r.db.WithContext(ctx).Model(&model.UserOrderStatistic{})
	if pageSize > 0 {
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Order("order_count DESC").Find(&stats).Error
	return stats, totalCount, err
}

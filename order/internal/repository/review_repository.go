package repository

import (
	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReviewRepository defines the interface for review repository operations
type ReviewRepository interface {
	CreateReview(review *model.Review) (*model.Review, error)
	GetReviewByID(id uuid.UUID) (*model.Review, error)
	GetReviewsByOrderID(orderID uuid.UUID) ([]model.Review, error)
	DeleteReview(id uuid.UUID) error
}

// reviewRepository implements the ReviewRepository interface
type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

// CreateReview creates a new review
func (r *reviewRepository) CreateReview(review *model.Review) (*model.Review, error) {
	if err := r.db.Create(review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

// GetReviewByID gets a review by ID
func (r *reviewRepository) GetReviewByID(id uuid.UUID) (*model.Review, error) {
	var review model.Review
	if err := r.db.Where("id = ?", id).First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

// GetReviewsByOrderID gets all reviews for an order
func (r *reviewRepository) GetReviewsByOrderID(orderID uuid.UUID) ([]model.Review, error) {
	var reviews []model.Review
	if err := r.db.Where("order_id = ?", orderID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// DeleteReview deletes a review
func (r *reviewRepository) DeleteReview(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Review{}).Error
}
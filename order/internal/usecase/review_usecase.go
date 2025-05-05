package usecase

import (
	"errors"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/google/uuid"
)

// reviewUseCase implements the ReviewUseCase interface
type reviewUseCase struct {
	reviewRepo repository.ReviewRepository
	orderRepo  repository.OrderRepository
}

// NewReviewUseCase creates a new review use case
func NewReviewUseCase(reviewRepo repository.ReviewRepository, orderRepo repository.OrderRepository) ReviewUseCase {
	return &reviewUseCase{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
	}
}

// CreateReview creates a new review
func (u *reviewUseCase) CreateReview(request model.CreateReviewRequest) (*model.Review, error) {
	// Verify that the order exists
	order, err := u.orderRepo.FindByID(request.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New(model.ErrOrderNotFound)
	}

	// Create the review
	review := &model.Review{
		OrderID:     request.OrderID,
		UserID:      request.UserID,
		Rating:      request.Rating,
		Description: request.Description,
	}

	return u.reviewRepo.CreateReview(review)
}

// GetReviewByID gets a review by ID
func (u *reviewUseCase) GetReviewByID(id uuid.UUID) (*model.Review, error) {
	return u.reviewRepo.GetReviewByID(id)
}

// GetReviewsByOrderID gets all reviews for an order
func (u *reviewUseCase) GetReviewsByOrderID(orderID uuid.UUID) ([]model.Review, error) {
	// Verify that the order exists
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New(model.ErrOrderNotFound)
	}

	return u.reviewRepo.GetReviewsByOrderID(orderID)
}

// DeleteReview deletes a review
func (u *reviewUseCase) DeleteReview(id uuid.UUID) error {
	// Verify that the review exists
	review, err := u.reviewRepo.GetReviewByID(id)
	if err != nil {
		return err
	}
	if review == nil {
		return errors.New(model.ErrReviewNotFound)
	}

	return u.reviewRepo.DeleteReview(id)
}

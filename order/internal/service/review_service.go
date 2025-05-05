package service

import (
	"net/http"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReviewService defines the interface for review service operations
type ReviewService interface {
	CreateReview(c *gin.Context)
	GetReviewByID(c *gin.Context)
	GetReviewsByOrderID(c *gin.Context)
	DeleteReview(c *gin.Context)
}

// reviewService implements the ReviewService interface
type reviewService struct {
	reviewRepo repository.ReviewRepository
}

// NewReviewService creates a new review service
func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

// CreateReview creates a new review
func (s *reviewService) CreateReview(c *gin.Context) {
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := &model.Review{
		OrderID:     req.OrderID,
		UserID:      req.UserID,
		Rating:      req.Rating,
		Description: req.Description,
	}

	createdReview, err := s.reviewRepo.CreateReview(review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.ReviewResponse{Review: *createdReview})
}

// GetReviewByID gets a review by ID
func (s *reviewService) GetReviewByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID format"})
		return
	}

	review, err := s.reviewRepo.GetReviewByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if review == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, model.ReviewResponse{Review: *review})
}

// GetReviewsByOrderID gets all reviews for an order
func (s *reviewService) GetReviewsByOrderID(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	reviews, err := s.reviewRepo.GetReviewsByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GetOrderReviewsResponse{Reviews: reviews})
}

// DeleteReview deletes a review
func (s *reviewService) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID format"})
		return
	}

	// Check if review exists
	review, err := s.reviewRepo.GetReviewByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if review == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	// Delete the review
	if err := s.reviewRepo.DeleteReview(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
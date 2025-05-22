package handler

import (
	"net/http"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReviewHandler handles HTTP requests for reviews
type ReviewHandler struct {
	reviewUseCase usecase.ReviewUseCase
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewUseCase usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{
		reviewUseCase: reviewUseCase,
	}
}

// CreateReview handles the request to create a new review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewUseCase.CreateReview(req)
	if err != nil {
		if err.Error() == model.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.ReviewResponse{Review: *review})
}

// GetReviewByID handles the request to get a review by ID
func (h *ReviewHandler) GetReviewByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID format"})
		return
	}

	review, err := h.reviewUseCase.GetReviewByID(id)
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

// GetReviewsByOrderID handles the request to get all reviews for an order
func (h *ReviewHandler) GetReviewsByOrderID(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	reviews, err := h.reviewUseCase.GetReviewsByOrderID(orderID)
	if err != nil {
		if err.Error() == model.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GetOrderReviewsResponse{Reviews: reviews})
}

// DeleteReview handles the request to delete a review
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID format"})
		return
	}

	if err := h.reviewUseCase.DeleteReview(id); err != nil {
		if err.Error() == model.ErrReviewNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// RegisterReviewRoutes registers the review routes
func RegisterReviewRoutes(router *gin.Engine, reviewHandler *ReviewHandler) {
	router.POST("/reviews", reviewHandler.CreateReview)
	router.GET("/reviews/:id", reviewHandler.GetReviewByID)
	router.GET("/orders/:orderId/reviews", reviewHandler.GetReviewsByOrderID)
	router.DELETE("/reviews/:id", reviewHandler.DeleteReview)
}

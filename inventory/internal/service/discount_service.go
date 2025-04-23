package service

import (
	"net/http"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscountService interface {
	CreateDiscount(c *gin.Context)
	GetDiscountByID(c *gin.Context)
	UpdateDiscount(c *gin.Context)
	DeleteDiscount(c *gin.Context)
	ListDiscounts(c *gin.Context)
	GetAllProductsWithPromotion(c *gin.Context)
	GetProductsByDiscountID(c *gin.Context)
}

type discountService struct {
	discountUseCase usecase.DiscountUseCase
}

func NewDiscountService(discountUseCase usecase.DiscountUseCase) DiscountService {
	return &discountService{
		discountUseCase: discountUseCase,
	}
}

func (s *discountService) CreateDiscount(c *gin.Context) {
	var req model.CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := s.discountUseCase.CreateDiscount(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mapDiscountToResponse(*discount)
	c.JSON(http.StatusCreated, response)
}

func (s *discountService) GetDiscountByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	discount, err := s.discountUseCase.GetDiscountByID(id)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mapDiscountToResponse(*discount)
	c.JSON(http.StatusOK, response)
}

func (s *discountService) UpdateDiscount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	var req model.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := s.discountUseCase.UpdateDiscount(id, req)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mapDiscountToResponse(*discount)
	c.JSON(http.StatusOK, response)
}

func (s *discountService) DeleteDiscount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	if err := s.discountUseCase.DeleteDiscount(id); err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted successfully"})
}

func (s *discountService) ListDiscounts(c *gin.Context) {
	discounts, err := s.discountUseCase.ListDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]model.DiscountResponse, len(discounts))
	for i, discount := range discounts {
		response[i] = mapDiscountToResponse(discount)
	}

	c.JSON(http.StatusOK, response)
}

func (s *discountService) GetAllProductsWithPromotion(c *gin.Context) {
	productsWithPromotion, err := s.discountUseCase.GetAllProductsWithPromotion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productsWithPromotion)
}

func (s *discountService) GetProductsByDiscountID(c *gin.Context) {
	discountIDStr := c.Param("id")
	discountID, err := uuid.Parse(discountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	products, err := s.discountUseCase.GetProductsByDiscountID(discountID)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Helper function to map a Discount entity to a DiscountResponse
func mapDiscountToResponse(discount model.Discount) model.DiscountResponse {
	return model.DiscountResponse{
		ID:                 discount.ID,
		Name:               discount.Name,
		Description:        discount.Description,
		DiscountPercentage: discount.DiscountPercentage,
		ApplicableProducts: discount.ApplicableProducts,
		StartDate:          discount.StartDate,
		EndDate:            discount.EndDate,
		IsActive:           discount.IsActive,
		CreatedAt:          discount.CreatedAt,
		UpdatedAt:          discount.UpdatedAt,
	}
}

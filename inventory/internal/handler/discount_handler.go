package handler

import (
	"net/http"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscountHandler struct {
	discountUseCase usecase.DiscountUseCase
}

func NewDiscountHandler(discountUseCase usecase.DiscountUseCase) *DiscountHandler {
	return &DiscountHandler{
		discountUseCase: discountUseCase,
	}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var request model.CreateDiscountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := h.discountUseCase.CreateDiscount(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, discount)
}

func (h *DiscountHandler) GetDiscountByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	discount, err := h.discountUseCase.GetDiscountByID(id)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discount)
}

func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var request model.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, err := h.discountUseCase.UpdateDiscount(id, request)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discount)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.discountUseCase.DeleteDiscount(id)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted successfully"})
}

func (h *DiscountHandler) ListDiscounts(c *gin.Context) {
	discounts, err := h.discountUseCase.ListDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discounts)
}

func (h *DiscountHandler) GetAllProductsWithPromotion(c *gin.Context) {
	productsWithPromotions, err := h.discountUseCase.GetAllProductsWithPromotion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productsWithPromotions)
}

func (h *DiscountHandler) GetProductsByDiscountID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	products, err := h.discountUseCase.GetProductsByDiscountID(id)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func RegisterDiscountRoutes(router *gin.Engine, discountHandler *DiscountHandler) {
	router.POST("/discounts", discountHandler.CreateDiscount)
	router.GET("/discounts/:id", discountHandler.GetDiscountByID)
	router.PATCH("/discounts/:id", discountHandler.UpdateDiscount)
	router.DELETE("/discounts/:id", discountHandler.DeleteDiscount)
	router.GET("/discounts", discountHandler.ListDiscounts)
	router.GET("/products/promotions", discountHandler.GetAllProductsWithPromotion)
	router.GET("/discounts/:id/products", discountHandler.GetProductsByDiscountID)
}

package handler

import (
	"github.com/baccala1010/e-commerce/inventory/internal/service"
	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	discountService service.DiscountService
}

func NewDiscountHandler(discountService service.DiscountService) *DiscountHandler {
	return &DiscountHandler{
		discountService: discountService,
	}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	h.discountService.CreateDiscount(c)
}

func (h *DiscountHandler) GetDiscountByID(c *gin.Context) {
	h.discountService.GetDiscountByID(c)
}

func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	h.discountService.UpdateDiscount(c)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	h.discountService.DeleteDiscount(c)
}

func (h *DiscountHandler) ListDiscounts(c *gin.Context) {
	h.discountService.ListDiscounts(c)
}

func (h *DiscountHandler) GetAllProductsWithPromotion(c *gin.Context) {
	h.discountService.GetAllProductsWithPromotion(c)
}

func (h *DiscountHandler) GetProductsByDiscountID(c *gin.Context) {
	h.discountService.GetProductsByDiscountID(c)
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

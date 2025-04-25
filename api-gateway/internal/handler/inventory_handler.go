package handler

import (
	"net/http"
	"strconv"

	"github.com/baccala1010/e-commerce/api-gateway/internal/adapter/grpc/client/inventory"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// InventoryHandler handles HTTP requests for the inventory service
type InventoryHandler struct {
	client *inventory.Client
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(client *inventory.Client) *InventoryHandler {
	return &InventoryHandler{
		client: client,
	}
}

// GetProduct handles GET /products/:id
func (h *InventoryHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	req := &pb.GetProductRequest{
		Id: id,
	}

	resp, err := h.client.GetProductByID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, resp.Product)
}

// ListProducts handles GET /products
func (h *InventoryHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID := c.Query("category_id")

	req := &pb.ListProductsRequest{
		Page:       int32(page),
		Limit:      int32(limit),
		CategoryId: categoryID,
	}

	resp, err := h.client.ListProducts(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to list products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
		"total":    resp.Total,
	})
}

// CreateProduct handles POST /products
func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req pb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, resp.Product)
}

// UpdateProduct handles PATCH /products/:id
func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	var req pb.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	resp, err := h.client.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to update product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, resp.Product)
}

// DeleteProduct handles DELETE /products/:id
func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	req := &pb.DeleteProductRequest{
		Id: id,
	}

	err = h.client.DeleteProduct(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to delete product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCategory handles GET /categories/:id
func (h *InventoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	req := &pb.GetCategoryRequest{
		Id: id,
	}

	resp, err := h.client.GetCategoryByID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}

	c.JSON(http.StatusOK, resp.Category)
}

// ListCategories handles GET /categories
func (h *InventoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.ListCategoriesRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}

	resp, err := h.client.ListCategories(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to list categories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": resp.Categories,
		"total":      resp.Total,
	})
}

// CreateCategory handles POST /categories
func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	var req pb.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, resp.Category)
}

// UpdateCategory handles PATCH /categories/:id
func (h *InventoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	var req pb.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	resp, err := h.client.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to update category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, resp.Category)
}

// DeleteCategory handles DELETE /categories/:id
func (h *InventoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	req := &pb.DeleteCategoryRequest{
		Id: id,
	}

	err = h.client.DeleteCategory(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to delete category: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Discount handlers

// CreateDiscount handles POST /discounts
func (h *InventoryHandler) CreateDiscount(c *gin.Context) {
	var req pb.CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateDiscount(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create discount: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create discount"})
		return
	}

	c.JSON(http.StatusCreated, resp.Discount)
}

// GetDiscountByID handles GET /discounts/:id
func (h *InventoryHandler) GetDiscountByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	req := &pb.GetDiscountRequest{
		Id: id,
	}

	resp, err := h.client.GetDiscountByID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get discount: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discount"})
		return
	}

	c.JSON(http.StatusOK, resp.Discount)
}

// UpdateDiscount handles PATCH /discounts/:id
func (h *InventoryHandler) UpdateDiscount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	var req pb.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	resp, err := h.client.UpdateDiscount(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to update discount: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update discount"})
		return
	}

	c.JSON(http.StatusOK, resp.Discount)
}

// DeleteDiscount handles DELETE /discounts/:id
func (h *InventoryHandler) DeleteDiscount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	req := &pb.DeleteDiscountRequest{
		Id: id,
	}

	err = h.client.DeleteDiscount(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to delete discount: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete discount"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllProductsWithPromotion handles GET /products/promotions
func (h *InventoryHandler) GetAllProductsWithPromotion(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.GetProductsWithPromotionRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}

	resp, err := h.client.GetAllProductsWithPromotion(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get products with promotions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products with promotions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
		"total":    resp.Total,
	})
}

// GetProductsByDiscountID handles GET /discounts/:id/products
func (h *InventoryHandler) GetProductsByDiscountID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount ID format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.GetProductsByDiscountIDRequest{
		DiscountId: id,
		Page:       int32(page),
		Limit:      int32(limit),
	}

	resp, err := h.client.GetProductsByDiscountID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get products by discount ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products by discount ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
		"total":    resp.Total,
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request model.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUseCase.CreateProduct(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	product, err := h.productUseCase.GetProductByID(id)
	if err != nil {
		if err.Error() == model.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var request model.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUseCase.UpdateProduct(id, request)
	if err != nil {
		if err.Error() == model.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.productUseCase.DeleteProduct(id)
	if err != nil {
		if err.Error() == model.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Parse query parameters
	var params repository.ListProductParams

	// Parse page and page size
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	params.Page = page
	params.PageSize = pageSize

	// Parse category ID
	categoryIDStr := c.Query("category_id")
	if categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err == nil {
			params.CategoryID = &categoryID
		}
	}

	// Parse price range
	minPriceStr := c.Query("min_price")
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			params.MinPrice = &minPrice
		}
	}

	maxPriceStr := c.Query("max_price")
	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	// Parse search term
	search := c.Query("search")
	if search != "" {
		params.Search = &search
	}

	products, total, err := h.productUseCase.ListProducts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func RegisterProductRoutes(router *gin.Engine, productHandler *ProductHandler) {
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products/:id", productHandler.GetProductByID)
	router.PATCH("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
	router.GET("/products", productHandler.ListProducts)
}

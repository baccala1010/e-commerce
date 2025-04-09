package handler

import (
	"github.com/baccala1010/e-commerce/inventory/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	h.productService.CreateProduct(c)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	h.productService.GetProductByID(c)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	h.productService.UpdateProduct(c)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	h.productService.DeleteProduct(c)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	h.productService.ListProducts(c)
}

func RegisterProductRoutes(router *gin.Engine, productHandler *ProductHandler) {
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products/:id", productHandler.GetProductByID)
	router.PATCH("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
	router.GET("/products", productHandler.ListProducts)
}

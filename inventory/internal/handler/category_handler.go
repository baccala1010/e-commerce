package handler

import (
	"github.com/baccala1010/e-commerce/inventory/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	h.categoryService.CreateCategory(c)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	h.categoryService.GetCategoryByID(c)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	h.categoryService.UpdateCategory(c)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	h.categoryService.DeleteCategory(c)
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	h.categoryService.ListCategories(c)
}

func RegisterCategoryRoutes(router *gin.Engine, categoryHandler *CategoryHandler) {
	router.POST("/categories", categoryHandler.CreateCategory)
	router.GET("/categories/:id", categoryHandler.GetCategoryByID)
	router.PATCH("/categories/:id", categoryHandler.UpdateCategory)
	router.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	router.GET("/categories", categoryHandler.ListCategories)
}

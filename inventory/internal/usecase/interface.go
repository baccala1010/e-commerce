package usecase

import (
	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/google/uuid"
)

// ProductUseCase defines the business logic for product operations
type ProductUseCase interface {
	CreateProduct(request model.CreateProductRequest) (*model.Product, error)
	GetProductByID(id uuid.UUID) (*model.Product, error)
	UpdateProduct(id uuid.UUID, request model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(id uuid.UUID) error
	ListProducts(params repository.ListProductParams) ([]model.Product, int64, error)
}

// CategoryUseCase defines the business logic for category operations
type CategoryUseCase interface {
	CreateCategory(request model.CreateCategoryRequest) (*model.Category, error)
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	UpdateCategory(id uuid.UUID, request model.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(id uuid.UUID) error
	ListCategories() ([]model.Category, error)
}

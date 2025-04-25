package frontend

import (
	"time"

	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"github.com/google/uuid"
)

// InventoryServiceClient defines the interface for the inventory gRPC client
type InventoryServiceClient interface {
	// Product methods
	CreateProduct(name, description string, price float64, stockLevel int, categoryID uuid.UUID) (*pb.Product, error)
	GetProductByID(productID uuid.UUID) (*pb.Product, error)
	UpdateProduct(productID uuid.UUID, name, description *string, price *float64, stockLevel *int, categoryID *uuid.UUID) (*pb.Product, error)
	DeleteProduct(productID uuid.UUID) error
	ListProducts(page, limit int, categoryID *uuid.UUID) ([]*pb.Product, int32, error)

	// Category methods
	CreateCategory(name, description string) (*pb.Category, error)
	GetCategoryByID(categoryID uuid.UUID) (*pb.Category, error)
	UpdateCategory(categoryID uuid.UUID, name, description *string) (*pb.Category, error)
	DeleteCategory(categoryID uuid.UUID) error
	ListCategories() ([]*pb.Category, int32, error)

	// Discount methods
	CreateDiscount(name, description string, discountPercentage float64, applicableProducts []uuid.UUID, startDate, endDate time.Time) (*pb.Discount, error)
	GetDiscountByID(discountID uuid.UUID) (*pb.Discount, error)
	UpdateDiscount(discountID uuid.UUID, name, description *string, discountPercentage *float64, applicableProducts []uuid.UUID, startDate, endDate *time.Time, isActive *bool) (*pb.Discount, error)
	DeleteDiscount(discountID uuid.UUID) error
	GetAllProductsWithPromotion() ([]*pb.Product, int32, error)
	GetProductsByDiscountID(discountID uuid.UUID) ([]*pb.Product, int32, error)
}

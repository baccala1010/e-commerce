package backoffice

import (
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
)

// InventoryServiceServer defines the interface for the inventory gRPC service
type InventoryServiceServer interface {
	pb.InventoryServiceServer
}

// Server represents the gRPC server for inventory service
type Server struct {
	pb.UnimplementedInventoryServiceServer
	productUseCase  usecase.ProductUseCase
	categoryUseCase usecase.CategoryUseCase
	discountUseCase usecase.DiscountUseCase
}

// NewServer creates a new inventory gRPC server
func NewServer(productUseCase usecase.ProductUseCase, categoryUseCase usecase.CategoryUseCase, discountUseCase usecase.DiscountUseCase) *Server {
	return &Server{
		productUseCase:  productUseCase,
		categoryUseCase: categoryUseCase,
		discountUseCase: discountUseCase,
	}
}

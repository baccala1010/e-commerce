package backoffice

import (
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
)

// OrderServiceServer defines the interface for the order gRPC service
type OrderServiceServer interface {
	pb.OrderServiceServer
}

// Server represents the gRPC server for order service
type Server struct {
	pb.UnimplementedOrderServiceServer
	orderUseCase usecase.OrderUseCase
}

// NewServer creates a new order gRPC server
func NewServer(orderUseCase usecase.OrderUseCase) *Server {
	return &Server{
		orderUseCase: orderUseCase,
	}
}

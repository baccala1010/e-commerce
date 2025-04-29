package usecase

import (
	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
)

// OrderUseCase represents the business logic interface for order operations
type OrderUseCase interface {
	CreateOrder(request model.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(id uuid.UUID) (*model.Order, error)
	UpdateOrderStatus(id uuid.UUID, request model.UpdateOrderStatusRequest) (*model.Order, error)
	ListUserOrders(userID uuid.UUID, page, pageSize int) ([]model.Order, int64, error)
}

type ReviewUseCase interface {
	CreateReview(request model.CreateReviewRequest) (*model.Review, error)
	GetReviewByID(id uuid.UUID) (*model.Review, error)
	GetReviewsByOrderID(orderID uuid.UUID) ([]model.Review, error)
	DeleteReview(id uuid.UUID) error
}

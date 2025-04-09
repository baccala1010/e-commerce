package usecase

import (
	"errors"
	"fmt"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/google/uuid"
)

type orderUseCase struct {
	orderRepo repository.OrderRepository
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(orderRepo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
	}
}

func (u *orderUseCase) CreateOrder(request model.CreateOrderRequest) (*model.Order, error) {
	// Create the order directly from the request
	order := &model.Order{
		UserID:        request.UserID,
		TotalAmount:   request.TotalAmount,
		ShippingName:  request.ShippingName,
		ShippingEmail: request.ShippingEmail,
		ShippingPhone: request.ShippingPhone,
		ShippingAddr:  request.ShippingAddr,
		Status:        model.OrderStatusPending,
	}

	// Create a payment record for the order
	payment := model.Payment{
		Amount: request.TotalAmount,
		Method: request.Payment.Method,
		Status: model.PaymentStatusPending,
	}

	// Associate payment with order
	order.Payment = payment

	if err := u.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return order, nil
}

func (u *orderUseCase) GetOrderByID(id uuid.UUID) (*model.Order, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding order: %w", err)
	}

	if order == nil {
		return nil, nil
	}

	return order, nil
}

func (u *orderUseCase) UpdateOrderStatus(id uuid.UUID, request model.UpdateOrderStatusRequest) (*model.Order, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding order: %w", err)
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, request.Status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, request.Status)
	}

	order.Status = request.Status

	if err := u.orderRepo.Update(order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	return order, nil
}

func (u *orderUseCase) ListUserOrders(userID uuid.UUID, page, pageSize int) ([]model.Order, int64, error) {
	return u.orderRepo.FindByUserID(userID, page, pageSize)
}

// isValidStatusTransition checks if the status transition is valid
func isValidStatusTransition(from, to model.OrderStatus) bool {
	validTransitions := map[model.OrderStatus][]model.OrderStatus{
		model.OrderStatusPending: {
			model.OrderStatusPaid,
			model.OrderStatusCancelled,
		},
		model.OrderStatusPaid: {
			model.OrderStatusShipped,
			model.OrderStatusCancelled,
		},
		model.OrderStatusShipped: {
			model.OrderStatusDelivered,
			model.OrderStatusCancelled,
		},
		model.OrderStatusDelivered: {},
		model.OrderStatusCancelled: {},
	}

	// Allow setting the same status
	if from == to {
		return true
	}

	// Check if transition is allowed
	for _, validStatus := range validTransitions[from] {
		if validStatus == to {
			return true
		}
	}

	return false
}

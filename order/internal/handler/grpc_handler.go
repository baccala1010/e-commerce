package handler

import (
	"context"
	"time"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	orderUseCase usecase.OrderUseCase
}

func NewGRPCHandler(orderUseCase usecase.OrderUseCase) *GRPCHandler {
	return &GRPCHandler{
		orderUseCase: orderUseCase,
	}
}

// Order methods
func (h *GRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	createReq := model.CreateOrderRequest{
		UserID:        userID,
		TotalAmount:   req.TotalAmount,
		ShippingName:  req.ShippingName,
		ShippingEmail: req.ShippingEmail,
		ShippingPhone: req.ShippingPhone,
		ShippingAddr:  req.ShippingAddress,
		Payment: model.PaymentDTO{
			Method: convertProtoPaymentMethodToModel(req.Payment.Method),
		},
	}

	order, err := h.orderUseCase.CreateOrder(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &pb.OrderResponse{
		Order: convertOrderToProto(order),
	}, nil
}

func (h *GRPCHandler) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	if order == nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	return &pb.OrderResponse{
		Order: convertOrderToProto(order),
	}, nil
}

func (h *GRPCHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	updateReq := model.UpdateOrderStatusRequest{
		Status: convertProtoOrderStatusToModel(req.Status),
	}

	order, err := h.orderUseCase.UpdateOrderStatus(orderID, updateReq)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		if err.Error() == "invalid status transition" {
			return nil, status.Errorf(codes.FailedPrecondition, "invalid status transition: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.OrderResponse{
		Order: convertOrderToProto(order),
	}, nil
}

func (h *GRPCHandler) ListUserOrders(ctx context.Context, req *pb.ListUserOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	orders, total, err := h.orderUseCase.ListUserOrders(userID, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user orders: %v", err)
	}

	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = convertOrderToProto(&order)
	}

	return &pb.ListOrdersResponse{
		Orders: protoOrders,
		Total:  int32(total),
	}, nil
}

// Payment methods - these are implemented directly in the handler since there's no corresponding usecase method
func (h *GRPCHandler) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.PaymentResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	// Get the order to ensure it exists and to get the payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	if order == nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	// In a real implementation, this would call a payment processor
	// For now, we'll just simulate a successful payment
	payment := order.Payment
	payment.Status = model.PaymentStatusSuccess
	payment.TransactionID = uuid.New().String()
	payment.PaymentDate = time.Now()

	// Update the order status to paid
	updateReq := model.UpdateOrderStatusRequest{
		Status: model.OrderStatusPaid,
	}
	_, err = h.orderUseCase.UpdateOrderStatus(orderID, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.PaymentResponse{
		Payment: convertPaymentToProto(&payment),
	}, nil
}

func (h *GRPCHandler) GetPaymentByID(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment ID: %v", err)
	}

	// In a real implementation, this would query the payment repository
	// For now, we'll return a not implemented error
	return nil, status.Errorf(codes.Unimplemented, "GetPaymentByID is not implemented")
}

func (h *GRPCHandler) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment ID: %v", err)
	}

	// In a real implementation, this would update the payment status in the repository
	// For now, we'll return a not implemented error
	return nil, status.Errorf(codes.Unimplemented, "UpdatePaymentStatus is not implemented")
}

// Helper functions to convert between model and proto
func convertOrderToProto(order *model.Order) *pb.Order {
	return &pb.Order{
		Id:              order.ID.String(),
		UserId:          order.UserID.String(),
		Status:          convertModelOrderStatusToProto(order.Status),
		TotalAmount:     order.TotalAmount,
		ShippingName:    order.ShippingName,
		ShippingEmail:   order.ShippingEmail,
		ShippingPhone:   order.ShippingPhone,
		ShippingAddress: order.ShippingAddr,
		Payment:         convertPaymentToProto(&order.Payment),
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}
}

func convertPaymentToProto(payment *model.Payment) *pb.Payment {
	return &pb.Payment{
		Id:            payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		Amount:        payment.Amount,
		Method:        convertModelPaymentMethodToProto(payment.Method),
		Status:        convertModelPaymentStatusToProto(payment.Status),
		TransactionId: payment.TransactionID,
		PaymentDate:   timestamppb.New(payment.PaymentDate),
		CreatedAt:     timestamppb.New(payment.CreatedAt),
		UpdatedAt:     timestamppb.New(payment.UpdatedAt),
	}
}

func convertModelOrderStatusToProto(status model.OrderStatus) pb.OrderStatus {
	switch status {
	case model.OrderStatusPending:
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case model.OrderStatusPaid:
		return pb.OrderStatus_ORDER_STATUS_PAID
	case model.OrderStatusShipped:
		return pb.OrderStatus_ORDER_STATUS_SHIPPED
	case model.OrderStatusDelivered:
		return pb.OrderStatus_ORDER_STATUS_DELIVERED
	case model.OrderStatusCancelled:
		return pb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func convertProtoOrderStatusToModel(status pb.OrderStatus) model.OrderStatus {
	switch status {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return model.OrderStatusPending
	case pb.OrderStatus_ORDER_STATUS_PAID:
		return model.OrderStatusPaid
	case pb.OrderStatus_ORDER_STATUS_SHIPPED:
		return model.OrderStatusShipped
	case pb.OrderStatus_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered
	case pb.OrderStatus_ORDER_STATUS_CANCELLED:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusPending
	}
}

func convertModelPaymentStatusToProto(status model.PaymentStatus) pb.PaymentStatus {
	switch status {
	case model.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case model.PaymentStatusSuccess:
		return pb.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case model.PaymentStatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	case model.PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func convertProtoPaymentStatusToModel(status pb.PaymentStatus) model.PaymentStatus {
	switch status {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return model.PaymentStatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		return model.PaymentStatusSuccess
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return model.PaymentStatusFailed
	case pb.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return model.PaymentStatusRefunded
	default:
		return model.PaymentStatusPending
	}
}

func convertModelPaymentMethodToProto(method model.PaymentMethod) pb.PaymentMethod {
	switch method {
	case model.PaymentMethodCreditCard:
		return pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodDebitCard:
		return pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD
	case model.PaymentMethodPaypal:
		return pb.PaymentMethod_PAYMENT_METHOD_PAYPAL
	case model.PaymentMethodBankWire:
		return pb.PaymentMethod_PAYMENT_METHOD_BANK_WIRE
	default:
		return pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func convertProtoPaymentMethodToModel(method pb.PaymentMethod) model.PaymentMethod {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD:
		return model.PaymentMethodDebitCard
	case pb.PaymentMethod_PAYMENT_METHOD_PAYPAL:
		return model.PaymentMethodPaypal
	case pb.PaymentMethod_PAYMENT_METHOD_BANK_WIRE:
		return model.PaymentMethodBankWire
	default:
		return model.PaymentMethodCreditCard
	}
}

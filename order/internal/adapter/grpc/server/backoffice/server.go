package backoffice

import (
	"context"
	"time"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Order methods
func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
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

	order, err := s.orderUseCase.CreateOrder(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &pb.OrderResponse{
		Order: convertOrderToProto(order),
	}, nil
}

func (s *Server) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	order, err := s.orderUseCase.GetOrderByID(orderID)
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

func (s *Server) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	updateReq := model.UpdateOrderStatusRequest{
		Status: convertProtoOrderStatusToModel(req.Status),
	}

	order, err := s.orderUseCase.UpdateOrderStatus(orderID, updateReq)
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

func (s *Server) ListUserOrders(ctx context.Context, req *pb.ListUserOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	orders, total, err := s.orderUseCase.ListUserOrders(userID, int(req.Page), int(req.Limit))
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

// Payment methods
func (s *Server) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.PaymentResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	// Get the order to ensure it exists and to get the payment details
	order, err := s.orderUseCase.GetOrderByID(orderID)
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
	_, err = s.orderUseCase.UpdateOrderStatus(orderID, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.PaymentResponse{
		Payment: convertPaymentToProto(&payment),
	}, nil
}

func (s *Server) GetPaymentByID(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment ID: %v", err)
	}

	// In a real implementation, this would query the payment repository
	// For now, we'll return a not implemented error
	return nil, status.Errorf(codes.Unimplemented, "GetPaymentByID is not implemented")
}

func (s *Server) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment ID: %v", err)
	}

	// In a real implementation, this would update the payment status in the repository
	// For now, we'll return a not implemented error
	return nil, status.Errorf(codes.Unimplemented, "UpdatePaymentStatus is not implemented")
}

// Review methods
func (s *Server) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	createReq := model.CreateReviewRequest{
		OrderID:     orderID,
		UserID:      userID,
		Rating:      convertProtoRatingToModel(req.Rating),
		Description: req.Description,
	}

	review, err := s.reviewUseCase.CreateReview(createReq)
	if err != nil {
		if err.Error() == model.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create review: %v", err)
	}

	return &pb.ReviewResponse{
		Review: convertReviewToProto(review),
	}, nil
}

func (s *Server) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.ReviewResponse, error) {
	reviewID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid review ID: %v", err)
	}

	review, err := s.reviewUseCase.GetReviewByID(reviewID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get review: %v", err)
	}

	if review == nil {
		return nil, status.Errorf(codes.NotFound, "review not found")
	}

	return &pb.ReviewResponse{
		Review: convertReviewToProto(review),
	}, nil
}

func (s *Server) GetOrderReviews(ctx context.Context, req *pb.GetOrderReviewsRequest) (*pb.GetOrderReviewsResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	reviews, err := s.reviewUseCase.GetReviewsByOrderID(orderID)
	if err != nil {
		if err.Error() == model.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order reviews: %v", err)
	}

	protoReviews := make([]*pb.Review, len(reviews))
	for i, review := range reviews {
		protoReviews[i] = convertReviewToProto(&review)
	}

	return &pb.GetOrderReviewsResponse{
		Reviews: protoReviews,
	}, nil
}

func (s *Server) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*emptypb.Empty, error) {
	reviewID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid review ID: %v", err)
	}

	if err := s.reviewUseCase.DeleteReview(reviewID); err != nil {
		if err.Error() == model.ErrReviewNotFound {
			return nil, status.Errorf(codes.NotFound, "review not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete review: %v", err)
	}

	return &emptypb.Empty{}, nil
}

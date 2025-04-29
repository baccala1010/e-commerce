package frontend

import (
	"github.com/baccala1010/e-commerce/order/pkg/pb"
)

// OrderServiceClient defines the interface for the order gRPC client
type OrderServiceClient interface {
	// Order methods
	CreateOrder(userID string, totalAmount float64, payment *pb.PaymentInfo, shippingInfo *ShippingInfo) (*pb.Order, error)
	GetOrderByID(orderID string) (*pb.Order, error)
	UpdateOrderStatus(orderID string, status pb.OrderStatus) (*pb.Order, error)
	ListUserOrders(userID string, page, limit int32) ([]*pb.Order, int32, error)

	// Payment methods
	ProcessPayment(orderID string, method pb.PaymentMethod) (*pb.Payment, error)
	GetPaymentByID(paymentID string) (*pb.Payment, error)
	UpdatePaymentStatus(paymentID string, status pb.PaymentStatus, transactionID string) (*pb.Payment, error)

	// Review methods
	CreateReview(orderID string, rating int32, comment string) (*pb.Review, error)
	GetReviewByID(reviewID string) (*pb.Review, error)
	GetReviewsByOrderID(orderID string) ([]*pb.Review, error)
	DeleteReview(reviewID string) error
}

// ShippingInfo contains shipping details for an order
type ShippingInfo struct {
	Name    string
	Email   string
	Phone   string
	Address string
}

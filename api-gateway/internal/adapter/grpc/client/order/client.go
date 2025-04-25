package order

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/baccala1010/e-commerce/api-gateway/pkg/grpcconn"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Client is a gRPC client for the order service
type Client struct {
	conn   *grpc.ClientConn
	client pb.OrderServiceClient
}

// NewClient creates a new order service client
func NewClient(ctx context.Context, connManager *grpcconn.ConnectionManager, cfg *config.Config) (*Client, error) {
	conn, err := connManager.GetConnection(ctx, "order", cfg.Services.Order)
	if err != nil {
		return nil, fmt.Errorf("failed to get order service connection: %w", err)
	}

	client := pb.NewOrderServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	logrus.Infof("Calling order service CreateOrder: %+v", req)
	return c.client.CreateOrder(ctx, req)
}

// GetOrderByID gets an order by ID
func (c *Client) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	logrus.Infof("Calling order service GetOrderByID: %+v", req)
	return c.client.GetOrderByID(ctx, req)
}

// UpdateOrderStatus updates an order status
func (c *Client) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	logrus.Infof("Calling order service UpdateOrderStatus: %+v", req)
	return c.client.UpdateOrderStatus(ctx, req)
}

// ListUserOrders lists orders for a user
func (c *Client) ListUserOrders(ctx context.Context, req *pb.ListUserOrdersRequest) (*pb.ListOrdersResponse, error) {
	logrus.Infof("Calling order service ListUserOrders: %+v", req)
	return c.client.ListUserOrders(ctx, req)
}

// ProcessPayment processes a payment for an order
func (c *Client) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.PaymentResponse, error) {
	logrus.Infof("Calling order service ProcessPayment: %+v", req)
	return c.client.ProcessPayment(ctx, req)
}

// GetPaymentByID gets a payment by ID
func (c *Client) GetPaymentByID(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	logrus.Infof("Calling order service GetPaymentByID: %+v", req)
	return c.client.GetPaymentByID(ctx, req)
}

// UpdatePaymentStatus updates a payment status
func (c *Client) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error) {
	logrus.Infof("Calling order service UpdatePaymentStatus: %+v", req)
	return c.client.UpdatePaymentStatus(ctx, req)
}

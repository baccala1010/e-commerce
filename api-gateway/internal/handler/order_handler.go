package handler

import (
	"net/http"
	"strconv"

	"github.com/baccala1010/e-commerce/api-gateway/internal/adapter/grpc/client/order"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// OrderHandler handles HTTP requests for the order service
type OrderHandler struct {
	client *order.Client
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(client *order.Client) *OrderHandler {
	return &OrderHandler{
		client: client,
	}
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	req := &pb.GetOrderRequest{
		Id: id,
	}

	resp, err := h.client.GetOrderByID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
		return
	}

	c.JSON(http.StatusOK, resp.Order)
}

// ListUserOrders handles GET /orders
func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	_, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &pb.ListUserOrdersRequest{
		UserId: userID,
		Page:   int32(page),
		Limit:  int32(limit),
	}

	resp, err := h.client.ListUserOrders(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to list user orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list user orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": resp.Orders,
		"total":  resp.Total,
	})
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req pb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, resp.Order)
}

// UpdateOrderStatus handles PATCH /orders/:id
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	var req pb.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	resp, err := h.client.UpdateOrderStatus(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to update order status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, resp.Order)
}

// ProcessPayment handles POST /orders/:id/payment
func (h *OrderHandler) ProcessPayment(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	_, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	var req pb.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.OrderId = orderID

	resp, err := h.client.ProcessPayment(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to process payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, resp.Payment)
}

// GetPayment handles GET /payments/:id
func (h *OrderHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	req := &pb.GetPaymentRequest{
		Id: id,
	}

	resp, err := h.client.GetPaymentByID(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to get payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment"})
		return
	}

	c.JSON(http.StatusOK, resp.Payment)
}

// UpdatePaymentStatus handles PATCH /payments/:id
func (h *OrderHandler) UpdatePaymentStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID is required"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	var req pb.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	resp, err := h.client.UpdatePaymentStatus(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to update payment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, resp.Payment)
}

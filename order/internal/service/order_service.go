package service

import (
	"net/http"
	"strconv"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	ListUserOrders(c *gin.Context)
}

type orderService struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	return &orderService{
		orderUseCase: orderUseCase,
	}
}

func (s *orderService) CreateOrder(c *gin.Context) {
	var request model.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.orderUseCase.CreateOrder(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (s *orderService) GetOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	order, err := s.orderUseCase.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *orderService) UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var request model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.orderUseCase.UpdateOrderStatus(id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *orderService) ListUserOrders(c *gin.Context) {
	userIDParam := c.Query("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse pagination params
	page, pageSize := getPaginationParams(c)

	orders, total, err := s.orderUseCase.ListUserOrders(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"size":   pageSize,
	})
}

// Helper function to parse pagination parameters
func getPaginationParams(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	pageParam := c.DefaultQuery("page", "1")
	if pageVal, err := strconv.Atoi(pageParam); err == nil && pageVal > 0 {
		page = pageVal
	}

	pageSizeParam := c.DefaultQuery("page_size", "10")
	if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil && pageSizeVal > 0 {
		pageSize = pageSizeVal
	}

	return page, pageSize
}

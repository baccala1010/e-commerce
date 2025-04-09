package main

import (
	"fmt"
	"net/http"

	"github.com/baccala1010/e-commerce/order/internal/app"
	"github.com/baccala1010/e-commerce/order/internal/config"
	"github.com/baccala1010/e-commerce/order/internal/database"
	"github.com/baccala1010/e-commerce/order/internal/handler"
	"github.com/baccala1010/e-commerce/order/internal/middleware"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up logging
	app.SetupLogging(cfg)

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize use cases
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Initialize handlers
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Register routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "service": cfg.Server.Name})
	})

	// Register order routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
			orders.GET("", orderHandler.ListUserOrders)
		}
	}

	// Start the HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("Starting %s server on %s", cfg.Server.Name, serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

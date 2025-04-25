package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baccala1010/e-commerce/api-gateway/internal/adapter/grpc/client/inventory"
	"github.com/baccala1010/e-commerce/api-gateway/internal/adapter/grpc/client/order"
	"github.com/baccala1010/e-commerce/api-gateway/internal/app"
	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/baccala1010/e-commerce/api-gateway/internal/handler"
	"github.com/baccala1010/e-commerce/api-gateway/internal/middleware"
	"github.com/baccala1010/e-commerce/api-gateway/pkg/grpcconn"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logrus.Infof("Received signal: %v", sig)
		cancel()
	}()

	// Load configuration
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up logging
	app.SetupLogging(cfg)

	// Initialize gRPC connection manager
	connManager := grpcconn.NewConnectionManager()
	defer connManager.Close()

	// Initialize gRPC clients
	inventoryClient, err := inventory.NewClient(ctx, connManager, cfg)
	if err != nil {
		logrus.Fatalf("Failed to create inventory client: %v", err)
	}

	orderClient, err := order.NewClient(ctx, connManager, cfg)
	if err != nil {
		logrus.Fatalf("Failed to create order client: %v", err)
	}

	// Initialize HTTP handlers
	inventoryHandler := handler.NewInventoryHandler(inventoryClient)
	orderHandler := handler.NewOrderHandler(orderClient)

	// Initialize service proxy for backward compatibility
	proxy := handler.NewServiceProxy(cfg)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "up",
			"service": "api-gateway",
		})
	})

	// Register inventory routes
	router.GET("/products", inventoryHandler.ListProducts)
	router.GET("/products/:id", inventoryHandler.GetProduct)
	router.POST("/products", inventoryHandler.CreateProduct)
	router.PATCH("/products/:id", inventoryHandler.UpdateProduct)
	router.DELETE("/products/:id", inventoryHandler.DeleteProduct)
	router.GET("/products/promotions", inventoryHandler.GetAllProductsWithPromotion)

	router.GET("/categories", inventoryHandler.ListCategories)
	router.GET("/categories/:id", inventoryHandler.GetCategory)
	router.POST("/categories", inventoryHandler.CreateCategory)
	router.PATCH("/categories/:id", inventoryHandler.UpdateCategory)
	router.DELETE("/categories/:id", inventoryHandler.DeleteCategory)

	// Register discount routes
	router.GET("/discounts/:id", inventoryHandler.GetDiscountByID)
	router.POST("/discounts", inventoryHandler.CreateDiscount)
	router.PATCH("/discounts/:id", inventoryHandler.UpdateDiscount)
	router.DELETE("/discounts/:id", inventoryHandler.DeleteDiscount)
	router.GET("/discounts/:id/products", inventoryHandler.GetProductsByDiscountID)

	// Register order routes
	router.GET("/orders", orderHandler.ListUserOrders)
	router.GET("/orders/:id", orderHandler.GetOrder)
	router.POST("/orders", orderHandler.CreateOrder)
	router.PATCH("/orders/:id", orderHandler.UpdateOrderStatus)
	router.POST("/orders/:id/payment", orderHandler.ProcessPayment)

	router.GET("/payments/:id", orderHandler.GetPayment)
	router.PATCH("/payments/:id", orderHandler.UpdatePaymentStatus)

	// Register legacy proxy routes for backward compatibility
	inventoryGroup := router.Group("/inventory")
	inventoryGroup.Any("/*path", proxy.ProxyInventory())

	orderGroup := router.Group("/order")
	orderGroup.Any("/*path", proxy.ProxyOrder())

	// Start the HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("Starting %s server on %s", cfg.Server.Name, serverAddr)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	logrus.Info("Shutting down server...")

	// Create a deadline for server shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exiting")
}

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baccala1010/e-commerce/order/internal/adapter/grpc/server/backoffice"
	"github.com/baccala1010/e-commerce/order/internal/app"
	"github.com/baccala1010/e-commerce/order/internal/config"
	"github.com/baccala1010/e-commerce/order/internal/database"
	"github.com/baccala1010/e-commerce/order/internal/handler"
	"github.com/baccala1010/e-commerce/order/internal/middleware"
	"github.com/baccala1010/e-commerce/order/internal/repository"
	"github.com/baccala1010/e-commerce/order/internal/usecase"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize use cases
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo, orderRepo)

	// Initialize handlers
	orderHandler := handler.NewOrderHandler(orderUseCase)
	reviewHandler := handler.NewReviewHandler(reviewUseCase)

	// Create backoffice gRPC server instance
	backofficeServer := backoffice.NewServer(orderUseCase, reviewUseCase)

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

	// Register routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
			orders.GET("", orderHandler.ListUserOrders)
			orders.GET("/:orderId/reviews", reviewHandler.GetReviewsByOrderID)
		}

		reviews := v1.Group("/reviews")
		{
			reviews.POST("", reviewHandler.CreateReview)
			reviews.GET("/:id", reviewHandler.GetReviewByID)
			reviews.DELETE("/:id", reviewHandler.DeleteReview)
		}
	}

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the gRPC server
	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, backofficeServer)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	logrus.Infof("Starting gRPC server on %s", grpcAddr)
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logrus.Errorf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		logrus.Infof("Starting HTTP server on :%d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Failed to serve HTTP: %v", err)
		}
	}()

	// Wait for termination signal
	<-signalChan
	logrus.Info("Received termination signal, shutting down...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	logrus.Info("gRPC server stopped")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("HTTP server shutdown error: %v", err)
	}
	logrus.Info("HTTP server stopped")
}

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

	"github.com/baccala1010/e-commerce/inventory/internal/adapter/grpc/server/backoffice"
	"github.com/baccala1010/e-commerce/inventory/internal/app"
	"github.com/baccala1010/e-commerce/inventory/internal/cache"
	"github.com/baccala1010/e-commerce/inventory/internal/config"
	"github.com/baccala1010/e-commerce/inventory/internal/database"
	"github.com/baccala1010/e-commerce/inventory/internal/handler"
	"github.com/baccala1010/e-commerce/inventory/internal/middleware"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
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
	baseProductRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	discountRepo := repository.NewDiscountRepository(db)

	// Initialize cache
	productCache := cache.NewMemoryCache()

	// Create cached repository
	productRepo := repository.NewCachedProductRepository(baseProductRepo, productCache)

	// Initialize cache with data
	cachedRepo, ok := productRepo.(*repository.CachedProductRepository)
	if ok {
		if err := cachedRepo.RefreshCache(); err != nil {
			logrus.Warnf("Failed to initialize product cache: %v", err)
		}

		// Set up periodic refresh (every 12 hours)
		productCache.StartPeriodicRefresh(12*time.Hour, cachedRepo.RefreshCache)
	}

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	discountUseCase := usecase.NewDiscountUseCase(discountRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	discountHandler := handler.NewDiscountHandler(discountUseCase)
	// Create backoffice gRPC server instance
	backofficeServer := backoffice.NewServer(productUseCase, categoryUseCase, discountUseCase)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "service": cfg.Server.Name})
	})

	// Register API routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProductByID)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/promotions", discountHandler.GetAllProductsWithPromotion)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.GET("", categoryHandler.ListCategories)
		}

		// Discount routes
		discounts := v1.Group("/discounts")
		{
			discounts.POST("", discountHandler.CreateDiscount)
			discounts.GET("/:id", discountHandler.GetDiscountByID)
			discounts.PATCH("/:id", discountHandler.UpdateDiscount)
			discounts.DELETE("/:id", discountHandler.DeleteDiscount)
			discounts.GET("", discountHandler.ListDiscounts)
			discounts.GET("/:id/products", discountHandler.GetProductsByDiscountID)
		}
	}

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the gRPC server
	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, backofficeServer)

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

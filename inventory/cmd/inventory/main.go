package main

import (
	"fmt"
	"net/http"

	"github.com/baccala1010/e-commerce/inventory/internal/app"
	"github.com/baccala1010/e-commerce/inventory/internal/config"
	"github.com/baccala1010/e-commerce/inventory/internal/database"
	"github.com/baccala1010/e-commerce/inventory/internal/handler"
	"github.com/baccala1010/e-commerce/inventory/internal/middleware"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/baccala1010/e-commerce/inventory/internal/service"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
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
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	// Initialize services
	productService := service.NewProductService(productUseCase)
	categoryService := service.NewCategoryService(categoryUseCase)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

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
	}

	// Start the HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("Starting %s server on %s", cfg.Server.Name, serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

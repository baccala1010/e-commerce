package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baccala1010/e-commerce/api-gateway/internal/handler"

	"github.com/baccala1010/e-commerce/api-gateway/internal/app"
	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/baccala1010/e-commerce/api-gateway/internal/middleware"
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

	// Initialize service proxy
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
	router.GET("/products", proxy.ProxyInventory())
	router.GET("/products/:id", proxy.ProxyInventory())
	router.POST("/products", proxy.ProxyInventory())
	router.PATCH("/products/:id", proxy.ProxyInventory())
	router.DELETE("/products/:id", proxy.ProxyInventory())
	router.GET("/products/promotions", proxy.ProxyInventory())

	router.GET("/categories", proxy.ProxyInventory())
	router.GET("/categories/:id", proxy.ProxyInventory())
	router.POST("/categories", proxy.ProxyInventory())
	router.PATCH("/categories/:id", proxy.ProxyInventory())
	router.DELETE("/categories/:id", proxy.ProxyInventory())

	// Register discount routes
	router.GET("/discounts/:id", proxy.ProxyInventory())
	router.POST("/discounts", proxy.ProxyInventory())
	router.PATCH("/discounts/:id", proxy.ProxyInventory())
	router.DELETE("/discounts/:id", proxy.ProxyInventory())
	router.GET("/discounts/:id/products", proxy.ProxyInventory())

	// Register order routes
	router.GET("/orders", proxy.ProxyOrder())
	router.GET("/orders/:id", proxy.ProxyOrder())
	router.POST("/orders", proxy.ProxyOrder())
	router.PATCH("/orders/:id", proxy.ProxyOrder())
	router.POST("/orders/:id/payment", proxy.ProxyOrder())
	router.GET("/orders/:id/reviews", proxy.ProxyOrder())

	router.GET("/payments/:id", proxy.ProxyOrder())
	router.PATCH("/payments/:id", proxy.ProxyOrder())

	// Register review routes
	router.POST("/reviews", proxy.ProxyOrder())
	router.GET("/reviews/:id", proxy.ProxyOrder())
	router.DELETE("/reviews/:id", proxy.ProxyOrder())

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

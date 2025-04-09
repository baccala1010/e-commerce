package main

import (
	"fmt"

	"github.com/baccala1010/e-commerce/api-gateway/internal/app"
	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/baccala1010/e-commerce/api-gateway/internal/handler"
	"github.com/baccala1010/e-commerce/api-gateway/internal/middleware"
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

	// Register routes
	handler.RegisterProxyRoutes(router, proxy)

	// Start the HTTP server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("Starting %s server on %s", cfg.Server.Name, serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

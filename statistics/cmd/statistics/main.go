package main

import (
	"log"

	"github.com/baccala1010/e-commerce/statistics/internal/app"
	"github.com/baccala1010/e-commerce/statistics/internal/config"
)

func main() {
	// Initialize logging
	app.InitializeLogging()
	
	// Get configuration path
	configPath, err := config.GetPath("")
	if err != nil {
		log.Fatalf("Failed to get config path: %v", err)
	}
	
	// Create and run the application
	application, err := app.New(configPath)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
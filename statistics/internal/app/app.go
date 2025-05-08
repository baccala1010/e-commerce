package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/baccala1010/e-commerce/statistics/internal/adapter/grpc"
	"github.com/baccala1010/e-commerce/statistics/internal/adapter/kafka"
	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/database"
	"github.com/baccala1010/e-commerce/statistics/internal/handler"
	"github.com/baccala1010/e-commerce/statistics/internal/repository"
	"github.com/baccala1010/e-commerce/statistics/internal/usecase"
)

// App represents the statistics application
type App struct {
	cfg           *config.Config
	grpcServer    *grpc.Server
	eventProcessor *kafka.EventProcessor
}

// New creates a new statistics application
func New(configPath string) (*App, error) {
	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// Setup application context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	db, err := database.InitDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize use cases
	statisticsUsecase := usecase.NewStatisticsUsecase(userRepo, orderRepo)

	// Initialize gRPC handler
	statisticsHandler := handler.NewStatisticsHandler(statisticsUsecase)

	// Initialize gRPC server
	grpcServer, err := grpc.NewServer(cfg, statisticsHandler)
	if err != nil {
		return nil, err
	}

	// Initialize Kafka event processor
	eventProcessor, err := kafka.NewEventProcessor(cfg, userRepo, orderRepo, productRepo)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:           cfg,
		grpcServer:    grpcServer,
		eventProcessor: eventProcessor,
	}, nil
}

// Run starts the statistics application
func (a *App) Run() error {
	log.Println("Starting Statistics Service")

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start Kafka consumer
	if err := a.eventProcessor.Start(ctx); err != nil {
		return err
	}
	defer a.eventProcessor.Stop()

	// Start gRPC server (blocking)
	if err := a.grpcServer.Start(ctx); err != nil {
		return err
	}

	return nil
}
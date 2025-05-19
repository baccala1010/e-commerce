package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Shopify/sarama"
	grpcadapter "github.com/baccala1010/e-commerce/statistics/internal/adapter/grpc"
	"github.com/baccala1010/e-commerce/statistics/internal/adapter/kafka"
	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/database"
	"github.com/baccala1010/e-commerce/statistics/internal/handler"
	"github.com/baccala1010/e-commerce/statistics/internal/repository"
	"github.com/baccala1010/e-commerce/statistics/internal/usecase"
)

func Run() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config/config.yaml"
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	repo := repository.NewStatisticsRepository(db)
	uc := usecase.NewStatisticsUsecase(repo)
	handler := handler.NewStatisticsHandler(uc)
	// Start gRPC server
	port := strconv.Itoa(cfg.GRPC.Port)
	go func() {
		if err := grpcadapter.RunGRPCServer(handler, port); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	// Start Kafka consumer
	consumer := kafka.NewConsumer(repo)
	group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, nil)
	if err != nil {
		log.Fatalf("failed to create kafka consumer group: %v", err)
	}
	ctx := context.Background()
	go func() {
		for {
			if err := group.Consume(ctx, cfg.Kafka.Topics, consumer); err != nil {
				log.Printf("kafka consume error: %v", err)
			}
		}
	}()
	// Wait for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down statistics service...")
	_ = group.Close()
}

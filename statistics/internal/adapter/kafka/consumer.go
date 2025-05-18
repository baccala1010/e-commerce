package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	eventspb "github.com/baccala1010/e-commerce/order/pkg/pb/event
	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/model"
	"github.com/baccala1010/e-commerce/statistics/internal/repository"
	kafkawrapper "github.com/baccala1010/e-commerce/statistics/pkg/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
)

// EventProcessor handles consuming events from Kafka and processing them
type EventProcessor struct {
	consumer          *kafka.Consumer
	userRepo          repository.UserRepository
	orderRepo         repository.OrderRepository
	productRepo       repository.ProductRepository
	orderEventTopic   string
	productEventTopic string
	userEventTopic    string
}

// NewEventProcessor creates a new Kafka event processor
func NewEventProcessor(
	cfg *config.Config,
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
) (*EventProcessor, error) {
	// Create Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		GroupID:          cfg.Kafka.ConsumerGroupID,
		AutoOffsetReset:  cfg.Kafka.AutoOffsetReset,
	}, []string{
		cfg.Kafka.Topics.OrderEvents,
		cfg.Kafka.Topics.ProductEvents,
		cfg.Kafka.Topics.UserEvents,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &EventProcessor{
		consumer:          kafkaConsumer,
		userRepo:          userRepo,
		orderRepo:         orderRepo,
		productRepo:       productRepo,
		orderEventTopic:   cfg.Kafka.Topics.OrderEvents,
		productEventTopic: cfg.Kafka.Topics.ProductEvents,
		userEventTopic:    cfg.Kafka.Topics.UserEvents,
	}, nil
}

// Start begins processing events from Kafka
func (p *EventProcessor) Start(ctx context.Context) error {
	// Start the Kafka consumer
	if err := p.consumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start Kafka consumer: %w", err)
	}

	log.Println("Kafka event processor started")

	// Process messages from the consumer
	go func() {
		for msg := range p.consumer.Messages() {
			log.Printf("Received message from topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

			switch *msg.TopicPartition.Topic {
			case p.orderEventTopic:
				if err := p.processOrderEvent(ctx, msg.Value); err != nil {
					log.Printf("Failed to process order event: %v", err)
				}
			case p.productEventTopic:
				if err := p.processProductEvent(ctx, msg.Value); err != nil {
					log.Printf("Failed to process product event: %v", err)
				}
			case p.userEventTopic:
				if err := p.processUserEvent(ctx, msg.Value); err != nil {
					log.Printf("Failed to process user event: %v", err)
				}
			}
		}
	}()

	return nil
}

// processOrderEvent handles order events
func (p *EventProcessor) processOrderEvent(ctx context.Context, data []byte) error {
	var orderEvent eventspb.OrderEvent
	if err := proto.Unmarshal(data, &orderEvent); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	orderData := orderEvent.GetOrderData()
	if orderData == nil {
		return fmt.Errorf("order data is nil in event")
	}

	// Create model.Order from protobuf Order
	order := model.Order{
		ID:          orderData.Id,
		UserID:      orderEvent.UserId,
		TotalAmount: float64(orderData.TotalAmount),
		OrderStatus: orderData.Status.String(),
		CreatedAt:   time.Now(), // Ideally should come from the event
		UpdatedAt:   time.Now(),
	}

	// Store or update the order based on event type
	switch orderEvent.EventType {
	case eventspb.EventType_EVENT_TYPE_CREATED:
		log.Printf("Processing order create event: %s", order.ID)
		if err := p.orderRepo.Create(ctx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Ensure the user exists (might have been created by a user event, but add if not)
		user := model.User{
			ID:               orderEvent.UserId,
			RegistrationDate: time.Now(), // Might not be accurate, but ensures the user exists
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := p.userRepo.Create(ctx, user); err != nil {
			log.Printf("Warning: failed to ensure user exists: %v", err)
		}

	case eventspb.EventType_EVENT_TYPE_UPDATED:
		log.Printf("Processing order update event: %s", order.ID)
		if err := p.orderRepo.Update(ctx, order); err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}
	}

	return nil
}

// processProductEvent handles product events
func (p *EventProcessor) processProductEvent(ctx context.Context, data []byte) error {
	var productEvent eventspb.ProductEvent
	if err := proto.Unmarshal(data, &productEvent); err != nil {
		return fmt.Errorf("failed to unmarshal product event: %w", err)
	}

	productData := productEvent.GetProductData()
	if productData == nil {
		return fmt.Errorf("product data is nil in event")
	}

	// Create model.Product from protobuf Product
	product := model.Product{
		ID:         productData.Id,
		Name:       productData.Name,
		CategoryID: productData.CategoryId,
		Price:      float64(productData.Price),
		CreatedAt:  time.Now(), // Ideally should come from the event
		UpdatedAt:  time.Now(),
	}

	// Store or update the product based on event type
	switch productEvent.EventType {
	case eventspb.EventType_EVENT_TYPE_CREATED:
		log.Printf("Processing product create event: %s", product.ID)
		if err := p.productRepo.Create(ctx, product); err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	case eventspb.EventType_EVENT_TYPE_UPDATED:
		log.Printf("Processing product update event: %s", product.ID)
		if err := p.productRepo.Update(ctx, product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
	}

	return nil
}

// processUserEvent handles user events
func (p *EventProcessor) processUserEvent(ctx context.Context, data []byte) error {
	var userEvent eventspb.UserEvent
	if err := proto.Unmarshal(data, &userEvent); err != nil {
		return fmt.Errorf("failed to unmarshal user event: %w", err)
	}

	// Create model.User from protobuf UserEvent
	user := model.User{
		ID:               userEvent.UserId,
		Email:            userEvent.Email,
		Name:             userEvent.Name,
		RegistrationDate: userEvent.RegistrationDate.AsTime(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Store or update the user based on event type
	switch userEvent.EventType {
	case eventspb.EventType_EVENT_TYPE_CREATED:
		log.Printf("Processing user create event: %s", user.ID)
		if err := p.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	case eventspb.EventType_EVENT_TYPE_UPDATED:
		log.Printf("Processing user update event: %s", user.ID)
		if err := p.userRepo.Update(ctx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

// Stop stops the Kafka consumer
func (p *EventProcessor) Stop() {
	if p.consumer != nil {
		p.consumer.Close()
	}
}

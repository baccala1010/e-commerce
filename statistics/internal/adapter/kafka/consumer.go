package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/baccala1010/e-commerce/statistics/internal/repository"
	"github.com/google/uuid"
)

type OrderEvent struct {
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"` // created, updated, deleted
}

type InventoryEvent struct {
	ProductID  string `json:"product_id"`
	EventType  string `json:"event_type"` // product_created, product_updated, product_deleted
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
}

type Consumer struct {
	repo repository.StatisticsRepository
}

func NewConsumer(repo repository.StatisticsRepository) *Consumer {
	return &Consumer{repo: repo}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))

		// Check which topic the message is from
		switch msg.Topic {
		case "order-events":
			c.handleOrderEvent(sess, msg)
		case "inventory-events":
			c.handleInventoryEvent(sess, msg)
		default:
			log.Printf("Unknown topic: %s", msg.Topic)
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

func (c *Consumer) handleOrderEvent(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var event OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal order event: %v", err)
		sess.MarkMessage(msg, "")
		return
	}

	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		log.Printf("Invalid user_id in order event: %v", err)
		sess.MarkMessage(msg, "")
		return
	}

	ctx := context.Background()

	switch event.EventType {
	case "created":
		log.Printf("Processing order created event for user %s", event.UserID)
		if err := c.repo.IncrementOrderCount(ctx, userID); err != nil {
			log.Printf("Failed to increment order count: %v", err)
			return
		}
	case "deleted":
		log.Printf("Processing order deleted event for user %s", event.UserID)
		if err := c.repo.DecrementOrderCount(ctx, userID); err != nil {
			log.Printf("Failed to decrement order count: %v", err)
			return
		}
	default:
		log.Printf("Unknown order event type: %s", event.EventType)
	}

	sess.MarkMessage(msg, "")
}

func (c *Consumer) handleInventoryEvent(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var event InventoryEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal inventory event: %v", err)
		sess.MarkMessage(msg, "")
		return
	}

	log.Printf("Received inventory event: type=%s, product=%s", event.EventType, event.ProductID)

	sess.MarkMessage(msg, "")
}

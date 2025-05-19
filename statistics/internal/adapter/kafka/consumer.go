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
		var event OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			continue
		}
		userID, err := uuid.Parse(event.UserID)
		if err != nil {
			log.Printf("invalid user_id: %v", err)
			continue
		}
		switch event.EventType {
		case "created":
			_ = c.repo.IncrementOrderCount(context.Background(), userID)
		case "deleted":
			_ = c.repo.DecrementOrderCount(context.Background(), userID)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

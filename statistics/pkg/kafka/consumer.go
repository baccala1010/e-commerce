package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerConfig struct {
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

type Consumer struct {
	consumer *kafka.Consumer
	topics   []string
	msgChan  chan *kafka.Message
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config ConsumerConfig, topics []string) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.BootstrapServers,
		"group.id":           config.GroupID,
		"auto.offset.reset":  config.AutoOffsetReset,
		"enable.auto.commit": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		msgChan:  make(chan *kafka.Message, 100),
	}, nil
}

// Start begins consuming messages from the configured topics
func (c *Consumer) Start(ctx context.Context) error {
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	log.Printf("Kafka consumer started, subscribed to topics: %v", c.topics)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Kafka consumer context done, stopping consumer...")
				c.consumer.Close()
				close(c.msgChan)
				return
			default:
				msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
				if err == nil {
					c.msgChan <- msg
				} else if !err.(kafka.Error).IsTimeout() {
					log.Printf("Consumer error: %v", err)
				}
			}
		}
	}()

	return nil
}

// Messages returns a channel of Kafka messages
func (c *Consumer) Messages() <-chan *kafka.Message {
	return c.msgChan
}

// Close closes the Kafka consumer and message channel
func (c *Consumer) Close() {
	if c.consumer != nil {
		c.consumer.Close()
	}
	close(c.msgChan)
}

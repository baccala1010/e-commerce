package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer represents a Kafka producer
type Producer struct {
	producer  *kafka.Producer
	topicName string
}

// NewProducer creates a new Kafka producer
func NewProducer(bootstrapServers, topicName string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Start monitoring delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Message delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &Producer{
		producer:  p,
		topicName: topicName,
	}, nil
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(key string, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topicName, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		Timestamp:      time.Now(),
	}

	return p.producer.Produce(message, nil)
}

// Close closes the Kafka producer
func (p *Producer) Close() {
	p.producer.Flush(5000) // Wait for any outstanding messages to be delivered
	p.producer.Close()
}
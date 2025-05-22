package kafka

import (
	"github.com/Shopify/sarama"
)

func NewProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(brokers, config)
}

package kafka

import (
	"L0/config"
	"L0/internal/models"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(cfg *config.Kafka) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.BootstrapServers})
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, topic: cfg.Topic}, nil
}

func (p *Producer) Publish(order models.OrderJSON) error {
	ord, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshaling order: %v", err)
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          ord,
	}, nil)
}

package kafka

import (
	"L0/config"
	"L0/internal/models"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewConsumer(cfg *config.Kafka) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          "order_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer, topic: cfg.Topic}, nil
}

func (c *Consumer) Subscribe() (*models.OrderJSON, error) {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return nil, fmt.Errorf("error subscribing to topic: %v", err)
	}

	msg, err := c.consumer.ReadMessage(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %v", err)
	}

	var order models.OrderJSON
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return nil, fmt.Errorf("error unmarshaling message: %v", err)
	}

	return &order, nil
}

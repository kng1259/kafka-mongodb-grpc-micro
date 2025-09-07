package kafka

import (
	"context"
	"encoding/json"
	"log"

	"kmgm-consumer/models"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Consume(ctx context.Context, processFunc func(models.Product) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}

		log.Printf("Received message: %s", string(msg.Value))

		// Parse product from message
		var product models.Product
		err = json.Unmarshal(msg.Value, &product)
		if err != nil {
			log.Printf("Error unmarshaling product: %v", err)
			continue
		}

		// Process the product using the provided function
		if err := processFunc(product); err != nil {
			log.Printf("Error processing product: %v", err)
			continue
		}
	}
}

package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"kmgm-consumer/models"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, username string, password string, groupID string) *Consumer {
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     false,
		SASLMechanism: mechanism,
	}

	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6, // 10MB
			Dialer:   dialer,
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

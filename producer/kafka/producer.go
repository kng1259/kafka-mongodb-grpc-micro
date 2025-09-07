package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}, nil
}

func (p *Producer) ProduceMessage(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

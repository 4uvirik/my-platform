package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/4uvirik/my-platform/pkg/events"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) PublishOrderCreated(ctx context.Context, evt events.OrderCreated) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(evt.OrderID),
		Value: data,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

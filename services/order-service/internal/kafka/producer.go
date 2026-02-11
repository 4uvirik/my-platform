package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.com/4uvirik/my-platform/pkg/events"
)

type Producer struct {
	writer *kafka.Writer
	log    *slog.Logger
	topic  string
}

func NewProducer(brokers []string, topic string, log *slog.Logger) (*Producer, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &Producer{writer: w, log: log, topic: topic}, nil
}

func (p *Producer) PublishOrderCreated(ctx context.Context, evt events.OrderCreated) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(evt.OrderID),
		Value: b,
		Time:  time.Now(),
	})
	if err != nil {
		p.log.Error("kafka publish failed",
			"topic", p.topic,
			"order_id", evt.OrderID,
			"error", err,
		)
		return err
	}

	p.log.Info("kafka event published",
		"topic", p.topic,
		"order_id", evt.OrderID,
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

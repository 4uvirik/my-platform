package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"gitlab.com/4uvirik/my-platform/pkg/events"
)

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewConsumer(brokers []string, topic, groupID string, logger *slog.Logger) (*Consumer, error) {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		logger: logger,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var evt events.OrderCreated
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			c.logger.Error("failed to unmarshal event", "error", err)
			continue
		}

		c.logger.Info("order created event received",
			"order_id", evt.OrderID,
			"user_id", evt.UserID,
			"amount", evt.Amount,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

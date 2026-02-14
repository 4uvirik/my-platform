package kafka

import (
	"context"
	"encoding/json"
	"gitlab.com/4uvirik/my-platform/services/notification-service/internal/dedup"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"gitlab.com/4uvirik/my-platform/pkg/events"
)

type Consumer struct {
	reader *kafka.Reader
	dedup  *dedup.Service
	logger *slog.Logger
}

func NewConsumer(brokers []string, topic, groupID string, dedupSvc *dedup.Service, logger *slog.Logger) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader: r,
		dedup:  dedupSvc,
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
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		duplicated, err := c.dedup.Seen(ctx, "order_created", evt.OrderID)
		if err != nil {
			c.logger.Error("dedup check failed", "err", err)
			continue
		}

		if duplicated {
			c.logger.Warn("duplicate event skipped", "order_id", evt.OrderID)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		c.logger.Info("order created event received",
			"order_id", evt.OrderID,
			"user_id", evt.UserID,
			"amount", evt.Amount,
		)

		// TODO: here will be real notification logic (email, push, etc)

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit kafka message", "err", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

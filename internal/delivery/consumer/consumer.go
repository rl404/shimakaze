package consumer

import (
	"context"
	"encoding/json"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// Consumer contains functions for consumer.
type Consumer struct {
	service service.Service
	pubsub  pubsub.PubSub
	topic   string
}

// New to create new consumer.
func New(service service.Service, ps pubsub.PubSub, topic string) (*Consumer, error) {
	return &Consumer{
		service: service,
		pubsub:  ps,
		topic:   topic,
	}, nil
}

// Subscribe to start subscribing to topic.
func (c *Consumer) Subscribe(nrApp *newrelic.Application) error {
	c.pubsub.Use(log.PubSubMiddlewareWithLog(utils.GetLogger(0), log.PubSubMiddlewareConfig{Error: true}))
	c.pubsub.Use(log.PubSubMiddlewareWithLog(utils.GetLogger(1), log.PubSubMiddlewareConfig{
		Topic:   c.topic,
		Payload: true,
		Error:   true,
	}))

	return c.pubsub.Subscribe(context.Background(), c.topic, func() pubsub.HandlerFunc {
		return func(ctx context.Context, message []byte) {
			var msg entity.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				_ = errors.Wrap(ctx, err)
				return
			}

			tx := nrApp.StartTransaction("Consumer " + string(msg.Type))
			defer tx.End()

			ctx = newrelic.NewContext(ctx, tx)

			if err := c.service.ConsumeMessage(ctx, msg); err != nil {
				_ = errors.Wrap(ctx, err)
			}
		}
	}())
}

// Close to stop consumer connection.
func (c *Consumer) Close() error {
	return c.pubsub.Close()
}

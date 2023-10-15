package consumer

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/errors/stack"
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
	ps.Use(log.PubSubMiddlewareWithLog(utils.GetLogger(0), log.PubSubMiddlewareConfig{Error: true}))
	ps.Use(log.PubSubMiddlewareWithLog(utils.GetLogger(1), log.PubSubMiddlewareConfig{
		Topic:   topic,
		Payload: true,
		Error:   true,
	}))

	return &Consumer{
		service: service,
		pubsub:  ps,
		topic:   topic,
	}, nil
}

// Subscribe to subscribe.
func (c *Consumer) Subscribe() error {
	return c.pubsub.Subscribe(context.Background(), c.topic, func(ctx context.Context, message []byte) {
		var msg entity.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			stack.Wrap(ctx, err)
			return
		}

		if err := c.service.ConsumeMessage(ctx, msg); err != nil {
			stack.Wrap(ctx, err)
		}
	})
}

// Close to stop consumer connection.
func (c *Consumer) Close() error {
	return c.pubsub.Close()
}

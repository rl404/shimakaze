package consumer

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// Consumer contains functions for consumer.
type Consumer struct {
	service service.Service
	pubsub  pubsub.PubSub
	nrApp   *newrelic.Application
}

// New to create new consumer.
func New(service service.Service, ps pubsub.PubSub, nrApp *newrelic.Application) (*Consumer, error) {
	ps.Use(log.PubSubMiddlewareWithLog(utils.GetLogger(0), log.PubSubMiddlewareConfig{Error: true}))
	return &Consumer{
		service: service,
		pubsub:  ps,
		nrApp:   nrApp,
	}, nil
}

// Close to stop consumer connection.
func (c *Consumer) Close() error {
	return c.pubsub.Close()
}

package consumer

import (
	"context"
	"encoding/json"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
)

// SubscribeParseAgency to subscribe parse agency.
func (c *Consumer) SubscribeParseAgency() error {
	return c.pubsub.Subscribe(context.Background(), entity.TopicParseAgency, middlewareWithLog(entity.TopicParseAgency, c.subscribeParseAgency))
}

func (c *Consumer) subscribeParseAgency(ctx context.Context, message []byte) {
	var msg entity.ParseAgencyRequest
	if err := json.Unmarshal(message, &msg); err != nil {
		stack.Wrap(ctx, err)
		return
	}

	tx := c.nrApp.StartTransaction("Consumer parse agency")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.service.ConsumeParseAgency(ctx, msg); err != nil {
		stack.Wrap(ctx, err)
	}
}

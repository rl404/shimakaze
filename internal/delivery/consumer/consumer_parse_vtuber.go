package consumer

import (
	"context"
	"encoding/json"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
)

// SubscribeParseVtuber to subscribe parse vtuber.
func (c *Consumer) SubscribeParseVtuber() error {
	return c.pubsub.Subscribe(context.Background(), entity.TopicParseVtuber, middlewareWithLog(entity.TopicParseVtuber, c.subscribeParseVtuber))
}

func (c *Consumer) subscribeParseVtuber(ctx context.Context, message []byte) {
	var msg entity.ParseVtuberRequest
	if err := json.Unmarshal(message, &msg); err != nil {
		stack.Wrap(ctx, err)
		return
	}

	tx := c.nrApp.StartTransaction("Consumer parse vtuber")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.service.ConsumeParseVtuber(ctx, msg); err != nil {
		stack.Wrap(ctx, err)
	}
}

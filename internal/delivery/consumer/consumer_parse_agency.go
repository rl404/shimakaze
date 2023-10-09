package consumer

import (
	"context"
	"encoding/json"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	"github.com/rl404/shimakaze/internal/utils"
)

// SubscribeParseAgency to subscribe parse agency.
func (c *Consumer) SubscribeParseAgency() error {
	return c.pubsub.Subscribe(context.Background(), entity.TopicParseAgency, c.subscribeParseAgency())
}

func (c *Consumer) subscribeParseAgency() pubsub.HandlerFunc {
	return log.PubSubHandlerFuncWithLog(utils.GetLogger(1), func(ctx context.Context, message []byte) {
		var msg entity.ParseAgencyRequest
		if err := json.Unmarshal(message, &msg); err != nil {
			_ = errors.Wrap(ctx, err)
			return
		}

		tx := c.nrApp.StartTransaction("Consumer parse agency")
		defer tx.End()

		ctx = newrelic.NewContext(ctx, tx)

		if err := c.service.ConsumeParseAgency(ctx, msg); err != nil {
			_ = errors.Wrap(ctx, err)
		}
	}, log.PubSubMiddlewareConfig{
		Topic:   entity.TopicParseAgency,
		Payload: true,
		Error:   true,
	})
}

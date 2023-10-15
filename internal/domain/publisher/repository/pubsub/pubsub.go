package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

// Pubsub contains functions for pubsub.
type Pubsub struct {
	pubsub pubsub.PubSub
	topic  string
}

// New to create new pubsub.
func New(ps pubsub.PubSub, topic string) *Pubsub {
	return &Pubsub{
		pubsub: ps,
		topic:  topic,
	}
}

// PublishParseVtuber to publish parse vtuber.
func (p *Pubsub) PublishParseVtuber(ctx context.Context, id int64, forced bool) error {
	msg, err := json.Marshal(entity.Message{
		Type:   entity.TypeParseVtuber,
		ID:     id,
		Forced: forced,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, msg); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}

// PublishParseAgency to publish parse agency.
func (p *Pubsub) PublishParseAgency(ctx context.Context, id int64, forced bool) error {
	msg, err := json.Marshal(entity.Message{
		Type:   entity.TypeParseAgency,
		ID:     id,
		Forced: forced,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, msg); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}

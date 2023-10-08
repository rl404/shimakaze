package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
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
func (p *Pubsub) PublishParseVtuber(ctx context.Context, data entity.ParseVtuberRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	return nil
}

// PublishParseAgency to publish parse agency.
func (p *Pubsub) PublishParseAgency(ctx context.Context, data entity.ParseAgencyRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	return nil
}

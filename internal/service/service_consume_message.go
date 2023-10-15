package service

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

// ConsumeMessage to consume pubsub message.
func (s *service) ConsumeMessage(ctx context.Context, msg entity.Message) error {
	switch msg.Type {
	case entity.TypeParseVtuber:
		return stack.Wrap(ctx, s.consumeParseVtuber(ctx, msg.ID, msg.Forced))
	case entity.TypeParseAgency:
		return stack.Wrap(ctx, s.consumeParseAgency(ctx, msg.ID, msg.Forced))
	default:
		return stack.Wrap(ctx, errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseVtuber(ctx context.Context, id int64, forced bool) error {
	if !forced {
		isOld, _, err := s.vtuber.IsOld(ctx, id)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateVtuber(ctx, id); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseAgency(ctx context.Context, id int64, forced bool) error {
	if !forced {
		isOld, _, err := s.agency.IsOld(ctx, id)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateAgency(ctx, id); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

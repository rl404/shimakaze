package service

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseVtuber:
		return errors.Wrap(ctx, s.consumeParseVtuber(ctx, data.Data))
	case entity.TypeParseAgency:
		return errors.Wrap(ctx, s.consumeParseAgency(ctx, data.Data))
	default:
		return errors.Wrap(ctx, _errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseVtuber(ctx context.Context, data []byte) error {
	var req entity.ParseVtuberRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, _errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		isOld, _, err := s.vtuber.IsOld(ctx, req.ID)
		if err != nil {
			return errors.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateVtuber(ctx, req.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseAgency(ctx context.Context, data []byte) error {
	var req entity.ParseAgencyRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, _errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		isOld, _, err := s.agency.IsOld(ctx, req.ID)
		if err != nil {
			return errors.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateAgency(ctx, req.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}

package service

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
)

// ConsumeParseVtuber to consume parse vtuber message.
func (s *service) ConsumeParseVtuber(ctx context.Context, data entity.ParseVtuberRequest) error {
	if !data.Forced {
		isOld, _, err := s.vtuber.IsOld(ctx, data.ID)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateVtuber(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

// ConsumeParseAgency to consume parse agency message.
func (s *service) ConsumeParseAgency(ctx context.Context, data entity.ParseAgencyRequest) error {
	if !data.Forced {
		isOld, _, err := s.agency.IsOld(ctx, data.ID)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateAgency(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

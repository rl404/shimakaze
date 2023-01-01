package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseVtuber(ctx context.Context, data entity.ParseVtuberRequest) error
	PublishParseAgency(ctx context.Context, data entity.ParseAgencyRequest) error
}

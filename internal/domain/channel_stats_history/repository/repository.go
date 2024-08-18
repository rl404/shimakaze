package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
)

// Repository contains functions for channel-history domain.
type Repository interface {
	Create(ctx context.Context, data entity.ChannelStats) (int, error)
}

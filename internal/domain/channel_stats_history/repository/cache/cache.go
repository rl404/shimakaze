package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/repository"
)

// Cache contains functions for channel-stats-history cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new channel-stats-history cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Create to create channel stats.
func (c *Cache) Create(ctx context.Context, data entity.ChannelStats) (int, error) {
	return c.repo.Create(ctx, data)
}

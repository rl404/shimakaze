package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
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

// Get to get channel stats.
func (c *Cache) Get(ctx context.Context, req entity.GetRequest) (data []entity.ChannelStats, code int, err error) {
	key := utils.GetKey("channel-stats-history", req.VtuberID, req.StartDate, req.EndDate)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.Get(ctx, req)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return data, code, nil
}

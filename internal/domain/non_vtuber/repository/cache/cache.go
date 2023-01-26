package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cache contains functions for non-vtuber cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new non-vtuber cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetAllIDs to get all ids.
func (c *Cache) GetAllIDs(ctx context.Context) (data []int64, code int, err error) {
	key := utils.GetKey("non-vtuber", "ids")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAllIDs(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// Create to create new non-vtuber.
func (c *Cache) Create(ctx context.Context, id int64) (int, error) {
	code, err := c.repo.Create(ctx, id)
	if err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("non-vtuber", "ids")
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return code, nil
}

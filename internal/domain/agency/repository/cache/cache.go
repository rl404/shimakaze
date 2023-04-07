package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/domain/agency/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cache contains functions for agency cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new agency cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByID to get by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Agency, code int, err error) {
	key := utils.GetKey("agency", id)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetAllIDs to get all ids.
func (c *Cache) GetAllIDs(ctx context.Context) (data []int64, code int, err error) {
	return c.repo.GetAllIDs(ctx)
}

// GetAll to get all.
func (c *Cache) GetAll(ctx context.Context) (data []entity.Agency, code int, err error) {
	key := utils.GetKey("agency")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAll(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// UpdateByID to update by id.
func (c *Cache) UpdateByID(ctx context.Context, id int64, data entity.Agency) (int, error) {
	code, err := c.repo.UpdateByID(ctx, id, data)
	if err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("agency", id)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return code, nil
}

// GetOldIDs to get old ids.
func (c *Cache) GetOldIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldIDs(ctx)
}

// GetCount to get count.
func (c *Cache) GetCount(ctx context.Context) (data int, code int, err error) {
	key := utils.GetKey("agency", "stats", "count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetCount(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

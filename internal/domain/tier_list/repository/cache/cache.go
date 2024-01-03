package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/tier_list/entity"
	"github.com/rl404/shimakaze/internal/domain/tier_list/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cache contains functions for tier list cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new tier list cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Get to get list.
func (c *Cache) Get(ctx context.Context, data entity.GetRequest) ([]entity.TierList, int, int, error) {
	return c.repo.Get(ctx, data)
}

// GetByID to get by id.
func (c *Cache) GetByID(ctx context.Context, id string) (data *entity.TierList, code int, err error) {
	key := utils.GetKey("tier-list", id)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return data, code, nil
}

// UpsertByID to upsert by id.
func (c *Cache) UpsertByID(ctx context.Context, data entity.TierList) (*entity.TierList, int, error) {
	res, code, err := c.repo.UpsertByID(ctx, data)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	key := utils.GetKey("tier-list", data.ID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)

	}

	return res, code, nil
}

// DeleteByID to delete by id.
func (c *Cache) DeleteByID(ctx context.Context, id string, userID int64) (int, error) {
	if code, err := c.repo.DeleteByID(ctx, id, userID); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	key := utils.GetKey("tier-list", id)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return http.StatusOK, nil
}

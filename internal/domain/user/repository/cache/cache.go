package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/user/entity"
	"github.com/rl404/shimakaze/internal/domain/user/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cache contains functions for user cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new user cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByID to get data by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.User, code int, err error) {
	key := utils.GetKey("user", id)
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

// Upsert to upsert.
func (c *Cache) Upsert(ctx context.Context, data entity.User) (int, error) {
	key := utils.GetKey("user", data.ID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}
	return c.repo.Upsert(ctx, data)
}

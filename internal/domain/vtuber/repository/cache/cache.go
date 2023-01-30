package cache

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cache contains functions for vtuber cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new vtuber cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByID to get data by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Vtuber, code int, err error) {
	key := utils.GetKey("vtuber", id)
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
func (c *Cache) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetAllIDs(ctx)
}

// GetAllImages to get all images.
func (c *Cache) GetAllImages(ctx context.Context, shuffle bool, limit int) (data []entity.Vtuber, code int, err error) {
	key := utils.GetKey("vtuber", "images")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAllImages(ctx, shuffle, limit)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(data), func(i, j int) {
			data[i], data[j] = data[j], data[i]
		})
	}

	if limit > 0 && len(data) > limit {
		data = data[:limit]
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// UpdateByID to update by id.
func (c *Cache) UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error) {
	return c.repo.UpdateByID(ctx, id, data)
}

// IsOld to check if old data.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// GetOldIDs to get old ids.
func (c *Cache) GetOldIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldIDs(ctx)
}

// GetAllForFamilyTree to get all for family tree.
func (c *Cache) GetAllForFamilyTree(ctx context.Context) (data []entity.Vtuber, code int, err error) {
	key := utils.GetKey("vtuber", "tree", "family")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAllForFamilyTree(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetAllForAgencyTree to get all for agency tree.
func (c *Cache) GetAllForAgencyTree(ctx context.Context) (data []entity.Vtuber, code int, err error) {
	key := utils.GetKey("vtuber", "tree", "agency")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAllForAgencyTree(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

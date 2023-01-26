package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/domain/agency/repository"
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

// GetAllIDs to get all ids.
func (c *Cache) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetAllIDs(ctx)
}

// GetAll to get all.
func (c *Cache) GetAll(ctx context.Context) ([]entity.Agency, int, error) {
	return c.repo.GetAll(ctx)
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// UpdateByID to update by id.
func (c *Cache) UpdateByID(ctx context.Context, id int64, data entity.Agency) (int, error) {
	return c.repo.UpdateByID(ctx, id, data)
}

// GetOldIDs to get old ids.
func (c *Cache) GetOldIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldIDs(ctx)
}

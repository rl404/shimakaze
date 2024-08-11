package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/shimakaze/internal/domain/language/entity"
	"github.com/rl404/shimakaze/internal/domain/language/repository"
)

// Cache contains functions for language cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new language cache.
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
func (c *Cache) GetAll(ctx context.Context) ([]entity.Language, int, int, error) {
	return c.repo.GetAll(ctx)
}

// UpdateByID to update by id.
func (c *Cache) UpdateByID(ctx context.Context, id int64, data entity.Language) (int, error) {
	return c.repo.UpdateByID(ctx, id, data)
}

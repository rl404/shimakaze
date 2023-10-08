package cache

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	_errors "github.com/rl404/shimakaze/internal/errors"
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(data), func(i, j int) {
			data[i], data[j] = data[j], data[i]
		})
	}

	if limit > 0 && len(data) > limit {
		data = data[:limit]
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// UpdateByID to update by id.
func (c *Cache) UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error) {
	if code, err := c.repo.UpdateByID(ctx, id, data); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("vtuber", id)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)

	}

	return http.StatusOK, nil
}

// DeleteByID to delete by id.
func (c *Cache) DeleteByID(ctx context.Context, id int64) (int, error) {
	if code, err := c.repo.DeleteByID(ctx, id); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("vtuber", id)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)

	}

	return http.StatusOK, nil
}

// IsOld to check if old data.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// GetOldActiveIDs to get old active ids.
func (c *Cache) GetOldActiveIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldActiveIDs(ctx)
}

// GetOldRetiredIDs to get old ids.
func (c *Cache) GetOldRetiredIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldRetiredIDs(ctx)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

type getAllCache struct {
	Data  []entity.Vtuber
	Total int
}

// GetAll to get vtuber list.
func (c *Cache) GetAll(ctx context.Context, req entity.GetAllRequest) (_ []entity.Vtuber, _ int, code int, err error) {
	key := utils.GetKey("vtuber", utils.QueryToKey(req))

	var data getAllCache
	if c.cacher.Get(ctx, key, &data) == nil {
		return data.Data, data.Total, http.StatusOK, nil
	}

	data.Data, data.Total, code, err = c.repo.GetAll(ctx, req)
	if err != nil {
		return nil, 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data.Data, data.Total, code, nil
}

// GetCharacterDesigners to get character designers.
func (c *Cache) GetCharacterDesigners(ctx context.Context) (data []string, code int, err error) {
	key := utils.GetKey("vtuber", "character-designers")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetCharacterDesigners(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetCharacter2DModelers to get 2d modelers.
func (c *Cache) GetCharacter2DModelers(ctx context.Context) (data []string, code int, err error) {
	key := utils.GetKey("vtuber", "character-2d-modelers")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetCharacter2DModelers(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetCharacter3DModelers to get 3d modelers.
func (c *Cache) GetCharacter3DModelers(ctx context.Context) (data []string, code int, err error) {
	key := utils.GetKey("vtuber", "character-3d-modelers")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetCharacter3DModelers(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

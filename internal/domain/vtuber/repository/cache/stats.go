package cache

import (
	"context"
	"net/http"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// GetCount to get count.
func (c *Cache) GetCount(ctx context.Context) (data int, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "count")
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

// GetAverageActiveTime to get average active time.
func (c *Cache) GetAverageActiveTime(ctx context.Context) (data float64, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "average-active-time")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAverageActiveTime(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetStatusCount to get status count.
func (c *Cache) GetStatusCount(ctx context.Context) (data *entity.StatusCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "status-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetStatusCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetDebutRetireCountMonthly to get debut & retire count monthly.
func (c *Cache) GetDebutRetireCountMonthly(ctx context.Context) (data []entity.DebutRetireCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "debut-retire-count-monthly")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetDebutRetireCountMonthly(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetDebutRetireCountYearly to get debut & retire count yearly.
func (c *Cache) GetDebutRetireCountYearly(ctx context.Context) (data []entity.DebutRetireCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "debut-retire-count-yearly")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetDebutRetireCountYearly(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetModelCount to get 2d & 3d model count.
func (c *Cache) GetModelCount(ctx context.Context) (data *entity.ModelCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "model-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetModelCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

// GetInAgencyCount to get in agency count.
func (c *Cache) GetInAgencyCount(ctx context.Context) (data *entity.InAgencyCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "in-agency-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetInAgencyCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}

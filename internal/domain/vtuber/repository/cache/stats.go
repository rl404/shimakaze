package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
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
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetSubscriberCount to get subscriber count.
func (c *Cache) GetSubscriberCount(ctx context.Context, interval, max int) (data []entity.SubscriberCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "subscriber-count", interval, max)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetSubscriberCount(ctx, interval, max)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetDesignerCount to get character designer count.
func (c *Cache) GetDesignerCount(ctx context.Context, top int) (data []entity.DesignerCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "designer-count", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetDesignerCount(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// Get2DModelerCount to get character 2d modeler count.
func (c *Cache) Get2DModelerCount(ctx context.Context, top int) (data []entity.DesignerCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "2d-modeler-count", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.Get2DModelerCount(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// Get3DModelerCount to get character 3d modeler count.
func (c *Cache) Get3DModelerCount(ctx context.Context, top int) (data []entity.DesignerCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "3d-modeler-count", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.Get3DModelerCount(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetAverageVideoCount to get average video count.
func (c *Cache) GetAverageVideoCount(ctx context.Context) (data float64, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "average-video-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAverageVideoCount(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetAverageVideoDuration to get average video duration.
func (c *Cache) GetAverageVideoDuration(ctx context.Context) (data float64, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "average-video-duration")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAverageVideoDuration(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetVideoCountByDate to get video count by date.
func (c *Cache) GetVideoCountByDate(ctx context.Context, hourly, daily bool) (data []entity.VideoCountByDate, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "video-count-by-date", hourly, daily)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetVideoCountByDate(ctx, hourly, daily)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetVideoCount to get video count.
func (c *Cache) GetVideoCount(ctx context.Context, top int) (data []entity.VideoCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "video-count", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetVideoCount(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetVideoDuration to get video duration.
func (c *Cache) GetVideoDuration(ctx context.Context, top int) (data []entity.VideoDuration, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "video-duration", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetVideoDuration(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetBirthdayCount to get birthday count.
func (c *Cache) GetBirthdayCount(ctx context.Context) (data []entity.BirthdayCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "birthday-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetBirthdayCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetAverageHeight to get average height.
func (c *Cache) GetAverageHeight(ctx context.Context) (data float64, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "average-height")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAverageHeight(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetAverageWeight to get average height.
func (c *Cache) GetAverageWeight(ctx context.Context) (data float64, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "average-weight")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetAverageWeight(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetBloodTypeCount to get blood type count.
func (c *Cache) GetBloodTypeCount(ctx context.Context, top int) (data []entity.BloodTypeCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "blood-type-count", top)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetBloodTypeCount(ctx, top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetChannelTypeCount to get channel type count.
func (c *Cache) GetChannelTypeCount(ctx context.Context) (data []entity.ChannelTypeCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "channel-type-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetChannelTypeCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetGenderCount to get gender count.
func (c *Cache) GetGenderCount(ctx context.Context) (data []entity.GenderCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "gender-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetGenderCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

// GetZodiacCount to get zodiac count.
func (c *Cache) GetZodiacCount(ctx context.Context) (data []entity.ZodiacCount, code int, err error) {
	key := utils.GetKey("vtuber", "stats", "zodiac-count")
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetZodiacCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return data, code, nil
}

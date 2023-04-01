package client

import (
	"context"
	_errors "errors"
	"net/http"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/rl404/shimakaze/internal/domain/twitch/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// GetVideos to get videos.
func (c *Client) GetVideos(ctx context.Context, id string) ([]entity.Video, int, error) {
	if code, err := c.setToken(ctx); err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	var res []entity.Video
	var cursor string
	for {
		resp, err := c.client.GetVideos(&helix.VideosParams{
			UserID: id,
			Type:   "archive",
			First:  100,
			After:  cursor,
		})
		if err != nil {
			if resp == nil {
				return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
			}
			return nil, resp.StatusCode, errors.Wrap(ctx, _errors.New(resp.Error), _errors.New(resp.ErrorMessage))
		}

		for _, v := range resp.Data.Videos {
			res = append(res, entity.Video{
				ID:        v.ID,
				Title:     v.Title,
				URL:       v.URL,
				Image:     c.getVideoImage(v.ThumbnailURL),
				StartDate: c.getStartDate(v.CreatedAt),
				EndDate:   c.getEndDate(v.CreatedAt, v.Duration),
			})
		}

		if len(resp.Data.Videos) < 100 {
			break
		}

		cursor = resp.Data.Pagination.Cursor
	}

	return res, http.StatusOK, nil
}

func (c *Client) getVideoImage(url string) string {
	url = strings.ReplaceAll(url, "%{width}", "400")
	url = strings.ReplaceAll(url, "%{height}", "200")
	return url
}

func (c *Client) getStartDate(d string) *time.Time {
	t, err := time.Parse(time.RFC3339, d)
	if err != nil {
		return nil
	}
	return &t
}

func (c *Client) getEndDate(start string, dur string) *time.Time {
	startDate := c.getStartDate(start)
	if startDate == nil {
		return nil
	}

	duration, err := utils.ParseDuration(dur, true)
	if err != nil || duration == 0 {
		return nil
	}

	endDate := startDate.Add(duration)

	return &endDate
}

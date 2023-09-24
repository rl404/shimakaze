package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rl404/shimakaze/internal/domain/youtube/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

type getVideosByIDsResponse struct {
	Items []video `json:"items"`
}

type video struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string          `json:"title"`
		Thumbnails  videoThumbnails `json:"thumbnails"`
		PublishedAt *time.Time      `json:"publishedAt"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration string
	} `json:"contentDetails"`
	LiveStreamingDetails struct {
		ActualStartTime    *time.Time `json:"actualStartTime"`
		ScheduledStartTime *time.Time `json:"scheduledStartTime"`
	} `json:"liveStreamingDetails"`
}

type videoThumbnails struct {
	Default thumbnail `json:"default"`
	Medium  thumbnail `json:"medium"`
}

// GetVideosByIDs to get videos by ids.
func (c *Client) GetVideosByIDs(ctx context.Context, ids []string) ([]entity.Video, int, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/videos", c.host))

	q := url.Query()
	q.Add("part", "snippet,contentDetails,liveStreamingDetails")

	var res []entity.Video
	idI := 0
	for {
		maxI := idI + 50
		if len(ids[idI:]) == 0 {
			break
		}

		if len(ids[idI:]) < 50 {
			maxI = idI + len(ids[idI:])
		}

		q.Set("id", strings.Join(ids[idI:maxI], ","))
		url.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		var body getVideosByIDsResponse
		if err := json.Unmarshal(respBody, &body); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		for _, item := range body.Items {
			res = append(res, entity.Video{
				ID:        item.ID,
				Title:     item.Snippet.Title,
				Image:     c.getVideoImage(item.Snippet.Thumbnails),
				StartDate: c.getVideoStartDate(item),
				EndDate:   c.getVideoEndDate(item),
			})
		}

		idI = maxI
	}

	return res, http.StatusOK, nil
}

func (c *Client) getVideoImage(thumbnails videoThumbnails) string {
	if thumbnails.Medium.URL != "" {
		return thumbnails.Medium.URL
	}
	return thumbnails.Default.URL
}

func (c *Client) getVideoStartDate(video video) *time.Time {
	if video.LiveStreamingDetails.ActualStartTime != nil {
		return video.LiveStreamingDetails.ActualStartTime
	}

	if video.LiveStreamingDetails.ScheduledStartTime != nil {
		return video.LiveStreamingDetails.ScheduledStartTime
	}

	return video.Snippet.PublishedAt
}

func (c *Client) getVideoEndDate(video video) *time.Time {
	startDate := c.getVideoStartDate(video)
	if startDate == nil {
		return nil
	}

	duration, err := utils.ParseDuration(video.ContentDetails.Duration)
	if err != nil || duration == 0 {
		return nil
	}

	endDate := startDate.Add(duration)

	return &endDate
}

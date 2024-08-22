package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/niconico/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getVideosResponse struct {
	Data struct {
		Items []struct {
			Essential struct {
				ID           string            `json:"id"`
				Title        string            `json:"title"`
				RegisteredAt time.Time         `json:"registeredAt"`
				Thumbnail    getVideosThumnail `json:"thumbnail"`
				Duration     int               `json:"duration"`
			} `json:"essential"`
		} `json:"items"`
	} `json:"data"`
}

type getVideosThumnail struct {
	URL       string `json:"url"`
	MiddleURL string `json:"middleUrl"`
	LargeURL  string `json:"largeUrl"`
}

// GetVideos to get videos.
func (c *Client) GetVideos(ctx context.Context, id string) ([]entity.Video, int, error) {
	url, _ := url.Parse(fmt.Sprintf("https://nvapi.nicovideo.jp/v3/users/%s/videos", id))

	q := url.Query()
	q.Add("sortKey", "registeredAt")
	q.Add("sortOrder", "desc")
	q.Add("pageSize", "100")

	// Loop until max age.
	var res []entity.Video
	page := 1
	for {
		q.Set("page", strconv.Itoa(page))
		url.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
		if err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		req.Header.Add("X-Frontend-id", "6")

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, resp.StatusCode, stack.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		var body getVideosResponse
		if err := json.Unmarshal(respBody, &body); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		var done bool
		for i, item := range body.Data.Items {
			if item.Essential.RegisteredAt.Before(c.maxAge) {
				done = true
				break
			}

			res = append(res, entity.Video{
				ID:        item.Essential.ID,
				Title:     item.Essential.Title,
				Image:     c.getVideoImage(item.Essential.Thumbnail),
				StartDate: &body.Data.Items[i].Essential.RegisteredAt,
				EndDate:   c.getVideoEndDate(item.Essential.RegisteredAt, item.Essential.Duration),
				URL:       fmt.Sprintf("https://www.nicovideo.jp/watch/%s", item.Essential.ID),
			})
		}

		if len(body.Data.Items) < 100 || done {
			break
		}

		page++
	}

	return res, http.StatusOK, nil
}

func (c *Client) getVideoImage(thumbnail getVideosThumnail) string {
	if thumbnail.LargeURL != "" {
		return thumbnail.LargeURL
	}

	if thumbnail.MiddleURL != "" {
		return thumbnail.MiddleURL
	}

	return thumbnail.URL
}

func (c *Client) getVideoEndDate(startDate time.Time, dur int) *time.Time {
	if dur <= 0 {
		return nil
	}
	endDate := startDate.Add(time.Duration(dur) * time.Second)
	return &endDate
}

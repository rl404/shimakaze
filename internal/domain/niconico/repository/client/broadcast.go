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

type getBroadcastsResponse struct {
	Data struct {
		ProgramsList []struct {
			ID struct {
				Value string `json:"value"`
			} `json:"id"`
			Program struct {
				Title    string `json:"title"`
				Schedule struct {
					BeginTime struct {
						Seconds int `json:"seconds"`
					} `json:"beginTime"`
					EndTime struct {
						Seconds int `json:"seconds"`
					} `json:"endTime"`
				} `json:"schedule"`
			} `json:"program"`
			Thumbnail getBroadcastsThumbnail `json:"thumbnail"`
		} `json:"programsList"`
	} `json:"data"`
}

type getBroadcastsThumbnail struct {
	Screenshot struct {
		Large  string `json:"large"`
		Middle string `json:"middle"`
		Small  string `json:"small"`
	} `json:"screenshot"`
}

// GetBroadcasts to get live broadcasts.
func (c *Client) GetBroadcasts(ctx context.Context, id string) ([]entity.Video, int, error) {
	url, _ := url.Parse("https://live.nicovideo.jp/front/api/v1/user-broadcast-history")

	q := url.Query()
	q.Add("providerId", id)
	q.Add("providerType", "user")
	q.Add("limit", "100")

	// Loop until max age.
	var res []entity.Video
	offset := 0
	for {
		q.Set("offset", strconv.Itoa(offset))
		url.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
		if err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

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

		var body getBroadcastsResponse
		if err := json.Unmarshal(respBody, &body); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		var done bool
		for _, item := range body.Data.ProgramsList {
			startDate := time.Unix(int64(item.Program.Schedule.BeginTime.Seconds), 0)
			endDate := time.Unix(int64(item.Program.Schedule.EndTime.Seconds), 0)

			if startDate.Before(c.maxAge) {
				done = true
				break
			}

			res = append(res, entity.Video{
				ID:        item.ID.Value,
				Title:     item.Program.Title,
				Image:     c.getBroadcastImage(item.Thumbnail),
				StartDate: &startDate,
				EndDate:   &endDate,
				URL:       fmt.Sprintf("https://live.nicovideo.jp/watch/%s", item.ID.Value),
			})
		}

		if len(body.Data.ProgramsList) < 100 || done {
			break
		}

		offset += 100
	}

	return res, http.StatusOK, nil
}

func (c *Client) getBroadcastImage(thumbnail getBroadcastsThumbnail) string {
	if thumbnail.Screenshot.Large != "" {
		return thumbnail.Screenshot.Large
	}

	if thumbnail.Screenshot.Middle != "" {
		return thumbnail.Screenshot.Middle
	}

	return thumbnail.Screenshot.Small
}

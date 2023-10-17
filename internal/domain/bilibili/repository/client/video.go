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
	"strings"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/bilibili/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getVideosResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List struct {
			VList []struct {
				BVID    string `json:"bvid"`
				Title   string `json:"title"`
				PIC     string `json:"pic"`
				Created int64  `json:"created"`
				Length  string `json:"length"`
			} `json:"vlist"`
		} `json:"list"`
	} `json:"data"`
}

// GetVideos to get videos.
func (c *Client) GetVideos(ctx context.Context, id string) ([]entity.Video, int, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/x/space/arc/search", c.host))

	q := url.Query()
	q.Add("mid", id)
	q.Add("order", "pubdate")
	q.Add("ps", "50")
	q.Add("tid", "0")

	var res []entity.Video
	page := 1
	for {
		q.Set("pn", strconv.Itoa(page))
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

		var body getVideosResponse
		if err := json.Unmarshal(respBody, &body); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		if body.Code != 0 {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, fmt.Errorf("%d %s", body.Code, body.Message), errors.ErrInternalServer)
		}

		var done bool
		for _, v := range body.Data.List.VList {
			startDate := c.getStartDate(v.Created)

			if startDate == nil || startDate.Before(c.maxAge) {
				done = true
				break
			}

			res = append(res, entity.Video{
				ID:        v.BVID,
				Title:     v.Title,
				Image:     v.PIC + "@200h",
				StartDate: startDate,
				EndDate:   c.getEndDate(startDate, v.Length),
			})
		}

		if len(body.Data.List.VList) < 50 || done {
			break
		}

		page++
	}

	return res, http.StatusOK, nil
}

func (c *Client) getStartDate(unix int64) *time.Time {
	if unix == 0 {
		return nil
	}

	t := time.Unix(unix, 0)
	return &t
}

func (c *Client) getEndDate(startDate *time.Time, dur string) *time.Time {
	if startDate == nil || dur == "" {
		return nil
	}

	splitDur := strings.Split(dur, ":")
	if len(splitDur) != 2 {
		return nil
	}

	m, err := strconv.Atoi(splitDur[0])
	if err != nil {
		return nil
	}

	s, err := strconv.Atoi(splitDur[1])
	if err != nil {
		return nil
	}

	endDate := startDate.Add(time.Duration(m)*time.Minute + time.Duration(s)*time.Second)

	return &endDate
}

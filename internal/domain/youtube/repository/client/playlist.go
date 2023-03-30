package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rl404/shimakaze/internal/errors"
)

type getPlaylistItemsResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		ContentDetails struct {
			VideoID          string    `json:"videoId"`
			VideoPublishedAt time.Time `json:"videoPublishedAt"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// GetVideoIDsByChannelID to get video ids by channel id.
func (c *Client) GetVideoIDsByChannelID(ctx context.Context, channelID string) ([]string, int, error) {
	playlistID := c.channelIDToPlaylistID(channelID)
	if playlistID == "" {
		return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrChannelNotFound)
	}

	url, _ := url.Parse(fmt.Sprintf("%s/playlistItems", c.host))

	q := url.Query()
	q.Add("playlistId", playlistID)
	q.Add("part", "contentDetails")
	q.Add("maxResults", "50")

	// Loop until max age.
	var res []string
	var nextToken string
	for {
		q.Set("pageToken", nextToken)
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

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		var body getPlaylistItemsResponse
		if err := json.Unmarshal(respBody, &body); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		var done bool
		for _, item := range body.Items {
			if item.ContentDetails.VideoPublishedAt.Before(c.maxAge) {
				done = true
				break
			}

			res = append(res, item.ContentDetails.VideoID)
		}

		if len(body.Items) == 0 || done {
			break
		}

		nextToken = body.NextPageToken
	}

	return res, http.StatusOK, nil
}

func (c *Client) channelIDToPlaylistID(channelID string) string {
	if len(channelID) < 3 {
		return ""
	}
	return "UU" + channelID[2:]
}

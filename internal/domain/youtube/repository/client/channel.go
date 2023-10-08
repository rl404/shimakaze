package client

import (
	"context"
	"encoding/json"
	__errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/youtube/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
)

type getChannelsByIDsResponse struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title      string            `json:"title"`
			Thumbnails channelThumbnails `json:"thumbnails"`
		} `json:"snippet"`
		Statistics struct {
			SubscriberCount string `json:"subscriberCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type channelThumbnails struct {
	Default thumbnail `json:"default"`
	Medium  thumbnail `json:"medium"`
}

type thumbnail struct {
	URL string `json:"url"`
}

// GetChannelByID to get channel by id.
func (c *Client) GetChannelByID(ctx context.Context, id string) (*entity.Channel, int, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/channels", c.host))

	q := url.Query()
	q.Add("id", id)
	q.Add("part", "snippet,contentDetails,statistics")
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Wrap(ctx, __errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	var body getChannelsByIDsResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	for _, channel := range body.Items {
		return &entity.Channel{
			ID:         channel.ID,
			Name:       channel.Snippet.Title,
			Image:      c.getChannelImage(channel.Snippet.Thumbnails),
			Subscriber: c.getSubscriber(channel.Statistics.SubscriberCount),
		}, http.StatusOK, nil
	}

	// No need to wrap the error to prevent useless error log.
	return nil, http.StatusNotFound, _errors.ErrChannelNotFound
}

func (c *Client) getChannelImage(thumbnails channelThumbnails) string {
	if thumbnails.Medium.URL != "" {
		return thumbnails.Medium.URL
	}
	return thumbnails.Default.URL
}

func (c *Client) getSubscriber(str string) int {
	subs, _ := strconv.Atoi(str)
	return subs
}

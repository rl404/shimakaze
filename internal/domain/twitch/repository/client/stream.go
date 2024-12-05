package client

import (
	"context"
	_errors "errors"
	"net/http"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/twitch/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

// GetStream to get live stream.
func (c *Client) GetLiveStream(ctx context.Context, id string) (*entity.Video, int, error) {
	if code, err := c.setToken(ctx); err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	resp, err := c.client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{id},
		First:   1,
	})
	if err != nil {
		if resp == nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}
		return nil, resp.StatusCode, stack.Wrap(ctx, _errors.New(resp.Error), _errors.New(resp.ErrorMessage))
	}

	for _, v := range resp.Data.Streams {
		return &entity.Video{
			ID:    v.ID,
			Image: c.getStreamImage(v.ThumbnailURL),
		}, http.StatusOK, nil
	}

	return nil, http.StatusOK, nil
}

func (c *Client) getStreamImage(url string) string {
	url = strings.ReplaceAll(url, "{width}", "400")
	url = strings.ReplaceAll(url, "{height}", "200")
	return url
}

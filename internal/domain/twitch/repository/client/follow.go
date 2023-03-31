package client

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"github.com/rl404/shimakaze/internal/errors"
)

// GetFollowerCount to get follower count.
func (c *Client) GetFollowerCount(ctx context.Context, id string) (int, int, error) {
	if code, err := c.setToken(ctx); err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	resp, err := c.client.GetChannelFollows(&helix.GetChannelFollowsParams{
		BroadcasterID: id,
		First:         1,
	})
	if err != nil {
		if resp == nil {
			return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}
		return 0, resp.StatusCode, errors.Wrap(ctx, _errors.New(resp.Error), _errors.New(resp.ErrorMessage))
	}

	return resp.Data.Total, http.StatusOK, nil
}

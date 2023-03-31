package client

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

func (c *Client) setToken(ctx context.Context) (int, error) {
	key := utils.GetKey("twitch", "token")

	// From cache.
	var token string
	if c.cacher.Get(ctx, key, &token) == nil {
		c.client.SetAppAccessToken(token)
		return http.StatusOK, nil
	}

	// Request token.
	resp, err := c.client.RequestAppAccessToken([]string{})
	if err != nil {
		if resp == nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}
		return resp.StatusCode, errors.Wrap(ctx, _errors.New(resp.Error), _errors.New(resp.ErrorMessage))
	}

	c.client.SetAppAccessToken(resp.Data.AccessToken)

	// Save to cache.
	if err := c.cacher.Set(ctx, key, resp.Data.AccessToken, time.Duration(resp.Data.ExpiresIn)*time.Second); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return http.StatusOK, nil
}

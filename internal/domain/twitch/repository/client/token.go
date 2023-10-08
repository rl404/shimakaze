package client

import (
	"context"
	__errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors"
	_errors "github.com/rl404/shimakaze/internal/errors"
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
			return http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
		}
		return resp.StatusCode, errors.Wrap(ctx, __errors.New(resp.Error), __errors.New(resp.ErrorMessage))
	}

	c.client.SetAppAccessToken(resp.Data.AccessToken)

	// Save to cache.
	if err := c.cacher.Set(ctx, key, resp.Data.AccessToken, time.Duration(resp.Data.ExpiresIn)*time.Second); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalCache)
	}

	return http.StatusOK, nil
}

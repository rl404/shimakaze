package client

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"github.com/rl404/shimakaze/internal/domain/twitch/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

// GetUser to get user.
func (c *Client) GetUser(ctx context.Context, name string) (*entity.User, int, error) {
	if code, err := c.setToken(ctx); err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	resp, err := c.client.GetUsers(&helix.UsersParams{Logins: []string{name}})
	if err != nil {
		if resp == nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}
		return nil, resp.StatusCode, errors.Wrap(ctx, _errors.New(resp.Error), _errors.New(resp.ErrorMessage))
	}

	if len(resp.Data.Users) == 0 {
		return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrChannelNotFound)
	}

	return &entity.User{
		ID:    resp.Data.Users[0].ID,
		Name:  resp.Data.Users[0].DisplayName,
		Image: resp.Data.Users[0].ProfileImageURL,
	}, http.StatusOK, nil
}

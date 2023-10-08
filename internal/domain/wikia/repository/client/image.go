package client

import (
	"context"
	"io"
	"net/http"

	"github.com/rl404/fairy/errors"
	_errors "github.com/rl404/shimakaze/internal/errors"
)

// GetImage to get image.
func (c *Client) GetImage(ctx context.Context, path string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Wrap(ctx, err)
	}

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	return image, http.StatusOK, nil
}

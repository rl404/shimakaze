package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rl404/shimakaze/internal/errors"
)

type getFollowerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Follower int `json:"follower"`
	} `json:"data"`
}

// GetFollowerCount to get follower count.
func (c *Client) GetFollowerCount(ctx context.Context, id string) (int, int, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/x/relation/stat?vmid=%s", c.host, id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	var body getFollowerResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if body.Code != 0 {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, fmt.Errorf("%d %s", body.Code, body.Message))
	}

	return body.Data.Follower, http.StatusOK, nil
}

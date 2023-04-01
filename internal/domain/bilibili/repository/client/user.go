package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rl404/shimakaze/internal/domain/bilibili/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MID  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"data"`
}

// GetUser to get user.
func (c *Client) GetUser(ctx context.Context, id string) (*entity.User, int, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/x/space/acc/info?mid=%s", c.host, id))

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

	var body getUserResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if body.Code != 0 {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, fmt.Errorf("%d %s", body.Code, body.Message))
	}

	return &entity.User{
		ID:    strconv.Itoa(body.Data.MID),
		Name:  body.Data.Name,
		Image: body.Data.Face,
	}, http.StatusOK, nil
}

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
	"github.com/rl404/shimakaze/internal/domain/bilibili/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
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
	url, _ := url.Parse(fmt.Sprintf("%s/x/space/wbi/acc/info?", c.host))

	q := url.Query()
	q.Add("mid", id)
	q.Add("platform", "web")
	q.Add("web_location", "1550101")
	q.Add("w_rid", "10bb0e85f7ff0dc7f03d5761206eba46")
	q.Add("wts", "1685590042")
	q.Add("token", "")
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

	var body getUserResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	if body.Code != 0 {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, fmt.Errorf("%d %s", body.Code, body.Message), _errors.ErrInternalServer)
	}

	return &entity.User{
		ID:    strconv.Itoa(body.Data.MID),
		Name:  body.Data.Name,
		Image: body.Data.Face,
	}, http.StatusOK, nil
}

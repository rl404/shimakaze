package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/wikia/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getPageImageByIDResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID   int64  `json:"pageid"`
			Title    string `json:"title"`
			Original struct {
				Source string `json:"source"`
			} `json:"original"`
		} `json:"pages"`
	} `json:"query"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetPageImageByID to get page image by id.
func (c *Client) GetPageImageByID(ctx context.Context, id int64) (*entity.PageImage, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("prop", "pageimages")
	q.Add("piprop", "original")
	q.Add("pageids", strconv.FormatInt(id, 10))
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, stack.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	var body getPageImageByIDResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if body.Error.Info != "" {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, _errors.New(body.Error.Info))
	}

	data, ok := body.Query.Pages[strconv.FormatInt(id, 10)]
	if !ok {
		return nil, http.StatusNotFound, stack.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	if data.Title == "" {
		return nil, http.StatusNotFound, stack.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	return &entity.PageImage{
		ID:    data.PageID,
		Title: data.Title,
		Image: data.Original.Source,
	}, http.StatusOK, nil
}

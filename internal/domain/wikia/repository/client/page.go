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

	"github.com/rl404/shimakaze/internal/domain/wikia/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getByIDResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID    int64  `json:"pageid"`
			Title     string `json:"title"`
			Revisions []struct {
				Slots struct {
					Main struct {
						Data string `json:"*"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetPageByID to get page by id.
func (c *Client) GetPageByID(ctx context.Context, id int64) (*entity.Page, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("prop", "revisions")
	q.Add("rvprop", "content")
	q.Add("rvslots", "main")
	q.Add("pageids", strconv.FormatInt(id, 10))
	url.RawQuery = q.Encode()

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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	var body getByIDResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if body.Error.Info != "" {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, _errors.New(body.Error.Info))
	}

	data, ok := body.Query.Pages[strconv.FormatInt(id, 10)]
	if !ok {
		return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	if data.Title == "" || len(data.Revisions) == 0 {
		return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	return &entity.Page{
		ID:      data.PageID,
		Title:   data.Title,
		Content: data.Revisions[0].Slots.Main.Data,
	}, http.StatusOK, nil
}

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

	"github.com/rl404/shimakaze/internal/domain/wikia/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getPageCategoriesResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID     int64  `json:"pageid"`
			Title      string `json:"title"`
			Categories []struct {
				Title string `json:"title"`
			} `json:"categories"`
		} `json:"pages"`
	} `json:"query"`
	Continue struct {
		CLContinue string `json:"clcontinue"`
	} `json:"continue"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetPageCategories to get page categories.
func (c *Client) GetPageCategories(ctx context.Context, id int64, limit int, lastTitle string) ([]entity.PageCategory, string, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("action", "query")
	q.Add("format", "json")
	q.Add("prop", "categories")
	q.Add("pageids", strconv.FormatInt(id, 10))
	q.Add("cllimit", strconv.Itoa(limit))

	if lastTitle != "" {
		q.Add("clcontinue", lastTitle)
	}

	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	var body getPageCategoriesResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if body.Error.Info != "" {
		return nil, "", http.StatusBadRequest, errors.Wrap(ctx, _errors.New(body.Error.Info))
	}

	data, ok := body.Query.Pages[strconv.FormatInt(id, 10)]
	if !ok {
		return nil, "", http.StatusNotFound, errors.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	if data.Title == "" {
		return nil, "", http.StatusNotFound, errors.Wrap(ctx, errors.ErrWikiaPageNotFound)
	}

	categories := make([]entity.PageCategory, len(data.Categories))
	for i, p := range data.Categories {
		categories[i] = entity.PageCategory{
			Title: p.Title,
		}
	}

	return categories, body.Continue.CLContinue, http.StatusOK, nil
}

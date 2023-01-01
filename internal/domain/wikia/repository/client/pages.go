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

type getPagesResponse struct {
	Query struct {
		AllPages []struct {
			PageID int64  `json:"pageid"`
			Title  string `json:"title"`
		} `json:"allpages"`
	} `json:"query"`
	Continue struct {
		APContinue string `json:"apcontinue"`
	} `json:"continue"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetPages to get pages.
func (c *Client) GetPages(ctx context.Context, limit int, lastName string) ([]entity.Page, string, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("action", "query")
	q.Add("format", "json")
	q.Add("list", "allpages")
	q.Add("apfilterredir", "nonredirects")
	q.Add("aplimit", strconv.Itoa(limit))
	q.Add("apcontinue", lastName)
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, err, errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, err, errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, err, errors.ErrInternalServer)
	}

	var body getPagesResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, "", http.StatusInternalServerError, errors.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if body.Error.Info != "" {
		return nil, "", http.StatusBadRequest, errors.Wrap(ctx, _errors.New(body.Error.Info))
	}

	pages := make([]entity.Page, len(body.Query.AllPages))
	for i, p := range body.Query.AllPages {
		pages[i] = entity.Page{
			ID:    p.PageID,
			Title: p.Title,
		}
	}

	return pages, body.Continue.APContinue, http.StatusOK, nil
}

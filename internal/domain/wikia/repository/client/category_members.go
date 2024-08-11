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

type getCategoryMembersResponse struct {
	Query struct {
		AllPages []struct {
			PageID int64  `json:"pageid"`
			Title  string `json:"title"`
		} `json:"categorymembers"`
	} `json:"query"`
	Continue struct {
		CMContinue string `json:"cmcontinue"`
	} `json:"continue"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetCategoryMembers to get category members.
func (c *Client) GetCategoryMembers(ctx context.Context, title string, limit int, lastTitle string, isPage bool) ([]entity.CategoryMember, string, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("action", "query")
	q.Add("format", "json")
	q.Add("list", "categorymembers")
	q.Add("cmtitle", title)
	q.Add("cmlimit", strconv.Itoa(limit))
	q.Add("cmcontinue", lastTitle)
	if isPage {
		q.Add("cmtype", "page")
		q.Add("cmnamespace", "0")
	}
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", resp.StatusCode, stack.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	var body getCategoryMembersResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if body.Error.Info != "" {
		return nil, "", http.StatusBadRequest, stack.Wrap(ctx, _errors.New(body.Error.Info))
	}

	members := make([]entity.CategoryMember, len(body.Query.AllPages))
	for i, p := range body.Query.AllPages {
		members[i] = entity.CategoryMember{
			ID:    p.PageID,
			Title: p.Title,
		}
	}

	return members, body.Continue.CMContinue, http.StatusOK, nil
}

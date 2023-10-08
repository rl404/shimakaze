package client

import (
	"context"
	"encoding/json"
	__errors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rl404/fairy/errors"
	_errors "github.com/rl404/shimakaze/internal/errors"
)

type getImageInfoResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID    int64  `json:"pageid"`
			Title     string `json:"title"`
			ImageInfo []struct {
				URL string `json:"url"`
			} `json:"imageinfo"`
		} `json:"pages"`
	} `json:"query"`
	Error struct {
		Info string `json:"info"`
	} `json:"error"`
}

// GetImageInfo to get image info.
func (c *Client) GetImageInfo(ctx context.Context, name string) (string, int, error) {
	c.limiter.Take()

	url, _ := url.Parse(fmt.Sprintf("%s/api.php", c.host))

	q := url.Query()
	q.Add("format", "json")
	q.Add("action", "query")
	q.Add("prop", "imageinfo")
	q.Add("iiprop", "url")
	q.Add("titles", name)
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.Wrap(ctx, __errors.New(http.StatusText(resp.StatusCode)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	var body getImageInfoResponse
	if err := json.Unmarshal(respBody, &body); err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalServer)
	}

	if body.Error.Info != "" {
		return "", http.StatusBadRequest, errors.Wrap(ctx, __errors.New(body.Error.Info))
	}

	if _, ok := body.Query.Pages["-1"]; ok {
		return "", http.StatusNotFound, errors.Wrap(ctx, _errors.ErrWikiaPageNotFound)
	}

	if len(body.Query.Pages) == 0 {
		return "", http.StatusNotFound, errors.Wrap(ctx, _errors.ErrWikiaPageNotFound)
	}

	for _, v := range body.Query.Pages {
		if len(v.ImageInfo) == 0 {
			continue
		}

		imgURL := v.ImageInfo[0].URL

		if imgURL == "" {
			return "", http.StatusNotFound, errors.Wrap(ctx, _errors.ErrWikiaPageNotFound)
		}

		return imgURL, http.StatusOK, nil
	}

	return "", http.StatusNotFound, errors.Wrap(ctx, _errors.ErrWikiaPageNotFound)
}

package client

import (
	"context"
	_errors "errors"
	"net/http"
	_url "net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/rl404/shimakaze/internal/errors"
)

// GetChannelIDByURL to get channel id by url.
func (c *Client) GetChannelIDByURL(ctx context.Context, url string) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", resp.StatusCode, errors.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	rssURLRaw := doc.Find("link[title=RSS]").AttrOr("href", "")

	rssURL, err := _url.Parse(rssURLRaw)
	if err != nil {
		return "", http.StatusNotFound, errors.Wrap(ctx, errors.ErrChannelNotFound)
	}

	channelID := rssURL.Query().Get("channel_id")

	if channelID == "" {
		return "", http.StatusNotFound, errors.Wrap(ctx, errors.ErrChannelNotFound)
	}

	return channelID, http.StatusOK, nil
}

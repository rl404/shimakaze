package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/niconico/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

type getUserResponse struct {
	State struct {
		UserDetails struct {
			UserDetails struct {
				User struct {
					FollowerCount int    `json:"followerCount"`
					ID            int    `json:"id"`
					Nickname      string `json:"nickname"`
					Icons         struct {
						Small string `json:"small"`
						Large string `json:"large"`
					} `json:"icons"`
				} `json:"user"`
			} `json:"userDetails"`
		} `json:"userDetails"`
	} `json:"state"`
}

// GetUser to get user.
func (c *Client) GetUser(ctx context.Context, userURL string) (*entity.User, int, error) {
	url, _ := url.Parse(userURL)

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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, stack.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	data, err := c.parseData(ctx, doc)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return &entity.User{
		ID:         strconv.Itoa(data.State.UserDetails.UserDetails.User.ID),
		Name:       data.State.UserDetails.UserDetails.User.Nickname,
		Image:      data.State.UserDetails.UserDetails.User.Icons.Large,
		Subscriber: data.State.UserDetails.UserDetails.User.FollowerCount,
	}, http.StatusOK, nil
}

func (c *Client) parseData(ctx context.Context, doc *goquery.Document) (*getUserResponse, error) {
	dataStr, ok := doc.Find("div#js-initial-userpage-data").Attr("data-initial-data")
	if !ok {
		return nil, stack.Wrap(ctx, errors.ErrChannelNotFound)
	}

	var data getUserResponse
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, stack.Wrap(ctx, err)
	}

	return &data, nil
}

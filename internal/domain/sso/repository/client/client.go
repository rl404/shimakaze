package client

import (
	"context"
	"encoding/json"
	_errors "errors"
	"io"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/sso/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"golang.org/x/oauth2"
)

type client struct {
	http   *http.Client
	config *oauth2.Config
}

// New to create new sso client.
func New(clientID, clientSecret, redirectURL string) *client {
	return &client{
		http: &http.Client{
			Timeout:   10 * time.Second,
			Transport: newrelic.NewRoundTripper(http.DefaultTransport),
		},
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint: oauth2.Endpoint{
				// AuthURL:  "https://sso.rl404.com",
				// TokenURL: "https://sso.rl404.com/api/oauth/token",
				AuthURL:  "http://localhost:3000",
				TokenURL: "http://localhost:3000/api/oauth/token",
			},
		},
	}
}

// ExchangeCode to exchange oauth code.
func (c *client) ExchangeCode(ctx context.Context, code string) (string, int, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}
	return token.AccessToken, http.StatusOK, nil
}

type userResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetUser to get user data.
func (c *client) GetUser(ctx context.Context, accessToken string) (*entity.User, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:41001/user", nil)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	req.Header.Add("authorization", "bearer "+accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, stack.Wrap(ctx, _errors.New(http.StatusText(resp.StatusCode)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	var userResp userResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return &entity.User{
		ID: userResp.ID,
	}, http.StatusOK, nil
}

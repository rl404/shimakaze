package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nicklaw5/helix/v2"
	"github.com/rl404/fairy/cache"
)

// Client is twitch api client.
type Client struct {
	cacher cache.Cacher
	client *helix.Client
}

// New to create new twitch api client.
func New(cacher cache.Cacher, clientID, clientSecret string) *Client {
	client, _ := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: newrelic.NewRoundTripper(http.DefaultTransport),
		},
	})
	return &Client{
		cacher: cacher,
		client: client,
	}
}

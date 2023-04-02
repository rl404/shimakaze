package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Client is niconico api client.
type Client struct {
	http   *http.Client
	maxAge time.Time
}

// New to create new niconico api client.
func New(maxAge int) *Client {
	return &Client{
		http: &http.Client{
			Timeout:   10 * time.Second,
			Transport: newrelic.NewRoundTripper(http.DefaultTransport),
		},
		maxAge: time.Now().Add(time.Duration(maxAge*-24) * time.Hour),
	}
}

package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/limit"
)

// Client contains functions for wikia api client.
type Client struct {
	host    string
	http    *http.Client
	limiter limit.Limiter
}

// New to create new wikia api client.
func New() *Client {
	l, _ := limit.New(limit.Mutex, 1, time.Second)
	return &Client{
		host: "https://virtualyoutuber.fandom.com",
		http: &http.Client{
			Timeout:   10 * time.Second,
			Transport: newrelic.NewRoundTripper(http.DefaultTransport),
		},
		limiter: l,
	}
}

package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Client contains functions for youtube api client.
type Client struct {
	host   string
	http   *http.Client
	maxAge time.Time
}

// New to create new youtube client.
func New(key string, maxAge int) *Client {
	return &Client{
		host: "https://www.googleapis.com/youtube/v3",
		http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: newrelic.NewRoundTripper(&transportWithKey{
				key: key,
			}),
		},
		maxAge: time.Now().Add(time.Duration(maxAge*-24) * time.Hour),
	}
}

type transportWithKey struct {
	transport http.RoundTripper
	key       string
}

// RoundTrip is http roundtrip.
func (t *transportWithKey) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.transport == nil {
		t.transport = http.DefaultTransport
	}

	q := req.URL.Query()
	q.Add("key", t.key)
	req.URL.RawQuery = q.Encode()

	return t.transport.RoundTrip(req)
}

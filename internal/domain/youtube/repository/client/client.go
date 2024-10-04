package client

import (
	"math/rand"
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
func New(keys []string, maxAge int) *Client {
	return &Client{
		host: "https://www.googleapis.com/youtube/v3",
		http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: newrelic.NewRoundTripper(&transportWithKey{
				randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
				keys:       keys,
			}),
		},
		maxAge: time.Now().Add(time.Duration(maxAge*-24) * time.Hour),
	}
}

type transportWithKey struct {
	transport  http.RoundTripper
	randomizer *rand.Rand
	keys       []string
}

// RoundTrip is http roundtrip.
func (t *transportWithKey) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.transport == nil {
		t.transport = http.DefaultTransport
	}

	q := req.URL.Query()
	q.Add("key", t.keys[t.randomizer.Intn(2)])
	req.URL.RawQuery = q.Encode()

	return t.transport.RoundTrip(req)
}

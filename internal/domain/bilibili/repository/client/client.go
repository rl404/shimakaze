package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Client is bilibili api client.
type Client struct {
	host   string
	http   *http.Client
	maxAge time.Time
}

// New to create new bilibili api client.
func New(maxAge int) *Client {
	return &Client{
		host: "https://api.bilibili.com",
		http: &http.Client{
			Timeout:   10 * time.Second,
			Transport: newrelic.NewRoundTripper(&transportWithHeader{}),
		},
		maxAge: time.Now().Add(time.Duration(maxAge*-24) * time.Hour),
	}
}

type transportWithHeader struct {
	transport http.RoundTripper
}

// RoundTrip is http roundtrip.
func (t *transportWithHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.transport == nil {
		t.transport = http.DefaultTransport
	}

	req.Header.Add("host", "api.bilibili.com")
	req.Header.Add("accept", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/111.0")

	return t.transport.RoundTrip(req)
}

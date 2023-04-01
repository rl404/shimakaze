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
func New(cookie string, maxAge int) *Client {
	return &Client{
		host: "https://api.bilibili.com",
		http: &http.Client{
			Timeout: 10 * time.Second,
			Transport: newrelic.NewRoundTripper(&transportWithCookie{
				cookie: cookie,
			}),
		},
		maxAge: time.Now().Add(time.Duration(maxAge*-24) * time.Hour),
	}
}

type transportWithCookie struct {
	transport http.RoundTripper
	cookie    string
}

// RoundTrip is http roundtrip.
func (t *transportWithCookie) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.transport == nil {
		t.transport = http.DefaultTransport
	}

	req.Header.Add("host", "api.bilibili.com")
	req.Header.Add("accept", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/111.0")
	req.Header.Add("cookie", t.cookie)

	return t.transport.RoundTrip(req)
}

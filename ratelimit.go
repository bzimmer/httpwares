package httpwares

import (
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimitTransport restricts the rate of api calls
type RateLimitTransport struct {
	Limiter   *rate.Limiter
	Transport http.RoundTripper
}

// RoundTrip executes requests within the rate limit
func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Limiter != nil {
		if err := t.Limiter.Wait(req.Context()); err != nil {
			return nil, err
		}
	}
	return transport(t.Transport).RoundTrip(req)
}

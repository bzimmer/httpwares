package transport

import (
	"net/http"
	"time"
)

// SleepingTransport is useful for testing delay scenarios
type SleepingTransport struct {
	Duration  time.Duration
	Transport http.RoundTripper
}

// RoundTrip sleeps for the specified duration then invokes the delegated RoundTripper
func (t *SleepingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	wait := make(chan bool)
	go func() {
		time.Sleep(t.Duration)
		wait <- true
	}()
	ctx := req.Context()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-wait:
		return t.Transport.RoundTrip(req)
	}
}

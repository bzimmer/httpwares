package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	// HeaderRateLimit header
	HeaderRateLimit = "X-Ratelimit-Limit"
	// HeaderRateUsage header
	HeaderRateUsage = "X-Ratelimit-Usage"
)

// RateLimit .
// See http://developers.strava.com/docs/rate-limits/ as an example
type RateLimit struct {
	LimitWindow int `json:"limit_window"`
	LimitDaily  int `json:"limit_daily"`
	UsageWindow int `json:"usage_window"`
	UsageDaily  int `json:"usage_daily"`
}

func (r *RateLimit) String() string {
	return fmt.Sprintf(
		"LimitWindow (%d), LimitDaily (%d), UsageWindow (%d), UsageDaily (%d)",
		r.LimitWindow, r.LimitDaily, r.UsageWindow, r.UsageDaily)
}

// RateLimitTransport .
type RateLimitTransport struct {
	RateLimit *RateLimit
	Transport http.RoundTripper
}

// RoundTrip .
func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.RateLimit != nil && t.RateLimit.IsThrottled() {
		return nil, t.RateLimit.NewError()
	}

	res, err := t.Transport.RoundTrip(req)

	if res != nil {
		rl, err := parseRateLimit(res)
		if err != nil {
			return nil, err
		}
		t.RateLimit = rl
	}

	return res, err
}

// RateLimitError .
type RateLimitError struct {
	RateLimit *RateLimit
}

func (e *RateLimitError) Error() string {
	return "exceeded rate limit"
}

func newRateLimitError(rl *RateLimit) *RateLimitError {
	return &RateLimitError{
		RateLimit: rl,
	}
}

// PercentDaily .
func (r *RateLimit) PercentDaily() int {
	if r.LimitDaily == 0 {
		return 0
	}
	return int(float32(r.UsageDaily) / float32(r.LimitDaily) * 100)
}

// PercentWindow .
func (r *RateLimit) PercentWindow() int {
	if r.LimitWindow == 0 {
		return 0
	}
	return int(float32(r.UsageWindow) / float32(r.LimitWindow) * 100)
}

// IsThrottled .
func (r *RateLimit) IsThrottled() bool {
	return r.PercentDaily() >= 100.0 || r.PercentWindow() >= 100.0
}

// NewError .
func (r *RateLimit) NewError() *RateLimitError {
	return newRateLimitError(r)
}

// parseRateLimit parses the headers returned from an API call into
// a RateLimit struct
//
//   HTTP/1.1 200 OK
//   Content-Type: application/json; charset=utf-8
//   Date: Tue, 10 Oct 2020 20:11:01 GMT
//   X-Ratelimit-Limit: 600,30000
//   X-Ratelimit-Usage: 314,27536
func parseRateLimit(res *http.Response) (*RateLimit, error) {
	var rateLimit RateLimit
	if limit := res.Header.Get(HeaderRateLimit); limit != "" {
		limits := strings.Split(limit, ",")
		x, err := strconv.Atoi(limits[0])
		if err != nil {
			return nil, err
		}
		rateLimit.LimitWindow = x
		x, err = strconv.Atoi(limits[1])
		if err != nil {
			return nil, err
		}
		rateLimit.LimitDaily = x
	}
	if usage := res.Header.Get(HeaderRateUsage); usage != "" {
		usages := strings.Split(usage, ",")
		x, err := strconv.Atoi(usages[0])
		if err != nil {
			return nil, err
		}
		rateLimit.UsageWindow = x
		x, err = strconv.Atoi(usages[1])
		if err != nil {
			return nil, err
		}
		rateLimit.UsageDaily = x
	}
	return &rateLimit, nil
}

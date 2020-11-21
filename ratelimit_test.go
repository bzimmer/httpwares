package httpwares_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"
)

func Test_RateLimitSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client := http.Client{
		Transport: &httpwares.RateLimitTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "transport.json",
				Status:      http.StatusOK,
				ContentType: "application/json",
				Responder: func(res *http.Response) error {
					res.Header.Add(httpwares.HeaderRateLimit, "600,30000")
					res.Header.Add(httpwares.HeaderRateUsage, "314,27536")
					return nil
				},
			},
		},
	}
	res, err := client.Get("http://example.com")
	a.NoError(err)
	a.NotNil(res)
}

func Test_RateLimitFailure(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	testdata := &httpwares.TestDataTransport{
		Filename:    "exceeded_rate_limit.json",
		Status:      http.StatusTooManyRequests,
		ContentType: "application/json",
		Responder: func(res *http.Response) error {
			res.Header.Add(httpwares.HeaderRateLimit, "575,30000")
			res.Header.Add(httpwares.HeaderRateUsage, "601,30100")
			return nil
		},
	}
	ratelimit := &httpwares.RateLimitTransport{
		Transport: testdata,
		RateLimit: &httpwares.RateLimit{},
	}

	client := http.Client{
		Transport: ratelimit,
	}

	// call the first time to seed the client with the rate limit response
	res, err := client.Get("http://example.com")
	a.NoError(err)
	a.True(ratelimit.RateLimit.IsThrottled())
	a.Equal("LimitWindow (575), LimitDaily (30000), UsageWindow (601), UsageDaily (30100)", ratelimit.RateLimit.String())
	a.NotNil(res)

	// the second call will fail not with the Fault but a RateLimitError
	//  (wrapped by url.Error) which can be inspected and used to throttle
	res, err = client.Get("http://example.com")
	a.Nil(res)
	a.Error(err)
	er := err.(*url.Error).Unwrap()
	a.Error(er.(*httpwares.RateLimitError))
	a.Equal("exceeded rate limit", er.Error())
	r := (er.(*httpwares.RateLimitError)).RateLimit
	a.Equal(30000, r.LimitDaily)
	a.Equal(601, r.UsageWindow)
	a.Equal(104, r.PercentWindow())
}

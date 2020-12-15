package httpwares_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	limiter := rate.NewLimiter(rate.Every(1*time.Minute), 1)
	client := http.Client{
		Transport: &httpwares.RateLimitTransport{
			Limiter: limiter,
			Transport: &httpwares.TestDataTransport{
				Status:      http.StatusOK,
				Filename:    "transport.json",
				ContentType: "application/json",
			}},
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	a.NoError(err)
	a.NotNil(req)

	res, err := client.Do(req)
	a.NoError(err)
	a.NotNil(res)

	ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*250)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	a.NoError(err)
	a.NotNil(req)

	// rate: Wait(n=1) would exceed context deadline
	res, err = client.Do(req)
	a.Error(err)
	a.Nil(res)
}

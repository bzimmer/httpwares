package httpwares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/httpwares"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/transport.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/transport.json")
	})

	tests := []struct {
		name string
		dur  time.Duration
	}{
		{
			name: "zero",
			dur:  0 * time.Millisecond,
		},
		{
			name: "with a duration",
			dur:  250 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(mux)
			defer svr.Close()

			limiter := rate.NewLimiter(rate.Every(1*time.Minute), 1)
			client := http.Client{
				Transport: &httpwares.RateLimitTransport{
					Limiter: limiter,
				},
			}

			ctx := context.Background()
			if tt.dur > (0 * time.Millisecond) {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, time.Millisecond*250)
				defer cancel()
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, svr.URL+"/transport.json", nil)
			a.NoError(err)
			a.NotNil(req)

			res, err := client.Do(req)
			a.NoError(err)
			a.NotNil(res)
			a.NoError(res.Body.Close())
		})
	}
}

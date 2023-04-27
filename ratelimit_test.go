package httpwares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"github.com/bzimmer/httpwares"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/transport.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/transport.json")
	})

	svr := httptest.NewServer(mux)
	defer svr.Close()

	limiter := rate.NewLimiter(rate.Every(1*time.Minute), 1)
	client := http.Client{
		Transport: &httpwares.RateLimitTransport{
			Limiter:   limiter,
			Transport: http.DefaultTransport,
		},
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svr.URL+"/transport.json", nil)
	a.NoError(err)
	a.NotNil(req)

	res, err := client.Do(req)
	a.NoError(err)
	a.NotNil(req)
	a.NoError(res.Body.Close())

	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, svr.URL+"/transport.json", nil)
	a.NoError(err)
	a.NotNil(req)

	res, err = client.Do(req)
	a.Error(err)
	a.Nil(res)
}

package httpwares_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client := http.Client{
		Transport: &httpwares.SleepingTransport{
			Duration: time.Millisecond * 100,
			Transport: &httpwares.TestDataTransport{
				Status:      http.StatusOK,
				Filename:    "transport.json",
				ContentType: "application/json",
			}},
	}

	// timeout lt sleep => failure
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*25)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	a.NoError(err)
	a.NotNil(req)

	res, err := client.Do(req)
	a.Error(err)
	a.Nil(res)
	if res != nil {
		a.NoError(res.Body.Close())
	}

	// timeout gt sleep => success
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, time.Millisecond*250)
	defer cancel()

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
	a.NoError(err)
	a.NotNil(req)

	res, err = client.Do(req)
	a.NoError(err)
	a.NotNil(res)
	a.NoError(res.Body.Close())
}

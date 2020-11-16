package transport_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bzimmer/transport"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client := http.Client{
		Transport: &transport.SleepingTransport{
			Duration: time.Millisecond * 100,
			Transport: &transport.TestDataTransport{
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
}

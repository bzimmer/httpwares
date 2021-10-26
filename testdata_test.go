package httpwares_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/httpwares"
)

func TestRequesterResponder(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		name, err, filename, contents string
		status                        int
		requester                     httpwares.Requester
		responder                     httpwares.Responder
	}{
		{
			name: "empty",
		},
		{
			name:     "valid",
			filename: "transport.txt",
			contents: `The mountains are calling & I must go & I will work on while I can, studying incessantly.`,
		},
		{
			name:     "valid filename but different status code",
			filename: "transport.txt",
			contents: `The mountains are calling & I must go & I will work on while I can, studying incessantly.`,
			status:   http.StatusBadRequest,
		},
		{
			name: "requester error",
			err:  "foo",
			requester: func(req *http.Request) error {
				return errors.New("foo")
			},
		},
		{
			name: "responder error",
			err:  "bar",
			responder: func(res *http.Response) error {
				return errors.New("bar")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			status := http.StatusOK
			if tt.status > 0 {
				status = tt.status
			}
			p := &httpwares.TestDataTransport{
				Filename:    "transport.txt",
				Status:      status,
				ContentType: "text/plain",
				Requester:   tt.requester,
				Responder:   tt.responder,
			}
			ctx := context.Background()
			client := http.Client{Transport: p}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
			a.NoError(err)
			a.NotNil(req)

			res, err := client.Do(req)
			if tt.err == "" {
				a.NoError(err)
				a.NotNil(res)
				a.Equal(status, res.StatusCode)
				defer res.Body.Close()
				data, err = io.ReadAll(res.Body)
				a.NoError(err)
				a.NotNil(data)
				a.Contains(string(data), tt.contents)
			} else {
				a.Error(err)
				a.Nil(res)
				a.Contains(err.Error(), tt.err)
				if res != nil {
					a.NoError(res.Body.Close())
				}
			}
		})
	}
}

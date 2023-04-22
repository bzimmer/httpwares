package httpwares_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/httpwares"
)

func TestVerboseTransport(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name        string
		filename    string
		status      int
		contenttype string
		body        string
		err         string
		writer      io.Writer
		requester   httpwares.Requester
	}{
		{
			name:        "plain no writer",
			filename:    "transport.txt",
			status:      http.StatusOK,
			contenttype: "text/plain",
		},
		{
			name:        "plain with writer",
			filename:    "transport.txt",
			status:      http.StatusOK,
			contenttype: "text/plain",
			writer:      new(bytes.Buffer),
			body:        "The mountains are calling & I must go & I will work on while I can, studying incessantly.",
		},
		{
			name:        "json",
			filename:    "transport.json",
			status:      http.StatusOK,
			contenttype: "application/json; charset=utf-8",
			writer:      new(bytes.Buffer),
			body:        `"quote": "The mountains are calling & I must go & I will work on while I can, studying incessantly."`,
		},
		{
			name:        "empty",
			filename:    "",
			status:      http.StatusNoContent,
			contenttype: "application/foobar; charset=utf-16",
		},
		{
			name: "error",
			err:  "argh",
			requester: func(req *http.Request) error {
				return errors.New("argh")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			v := &httpwares.VerboseTransport{
				Writer: tt.writer,
				Transport: &httpwares.TestDataTransport{
					Filename:    tt.filename,
					ContentType: tt.contenttype,
					Status:      tt.status,
					Requester:   tt.requester,
				},
			}
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
			a.NoError(err)
			a.NotNil(req)
			client := http.Client{Transport: v}
			res, err := client.Do(req)
			if tt.err != "" {
				a.Error(err)
				a.Nil(res)
				a.Contains(err.Error(), tt.err)
				return
			}
			a.NoError(err)
			a.NotNil(res)
			defer res.Body.Close()
			if tt.writer != nil && tt.body != "" {
				a.Contains(tt.writer.(*bytes.Buffer).String(), tt.body)
			}
		})
	}
}

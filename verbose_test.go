package httpwares_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/httpwares"
)

func TestVerboseTransport(t *testing.T) {
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/transport.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/transport.txt")
	})
	mux.HandleFunc("/transport.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/transport.json")
	})
	mux.HandleFunc("/Nikon_D70.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/Nikon_D70.jpg")
	})

	for _, tt := range []struct {
		name     string
		filename string
		status   int
		body     string
		writer   io.Writer
	}{
		{
			name:     "plain no writer",
			filename: "transport.txt",
			status:   http.StatusOK,
		},
		{
			name:     "plain with writer",
			filename: "transport.txt",
			status:   http.StatusOK,
			writer:   new(bytes.Buffer),
			body:     "The mountains are calling & I must go & I will work on while I can, studying incessantly.",
		},
		{
			name:     "json",
			filename: "transport.json",
			status:   http.StatusOK,
			writer:   new(bytes.Buffer),
			body:     `"quote": "The mountains are calling & I must go & I will work on while I can, studying incessantly."`,
		},
		{
			name:     "json",
			filename: "Nikon_D70.jpg",
			status:   http.StatusOK,
			writer:   new(bytes.Buffer),
		},
		{
			name:     "empty",
			filename: "",
			status:   http.StatusNoContent,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(mux)
			defer svr.Close()

			client := http.Client{
				Transport: &httpwares.VerboseTransport{Writer: tt.writer},
			}

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, svr.URL+"/"+tt.filename, nil)
			a.NoError(err)
			a.NotNil(req)

			res, err := client.Do(req)
			a.NoError(err)
			a.NotNil(res)
			defer res.Body.Close()

			if tt.writer != nil {
				a.Contains(tt.writer.(*bytes.Buffer).String(), tt.body)
			}
		})
	}
}

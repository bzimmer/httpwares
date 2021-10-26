package httpwares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

// Requester is a callback for http requests
type Requester func(*http.Request) error

// Responder is a callback for http responses
type Responder func(*http.Response) error

// TestDataTransport simplifies mocking http round trips
type TestDataTransport struct {
	Status      int
	Filename    string
	ContentType string
	Requester   Requester
	Responder   Responder
}

// RoundTrip responds to requests with mocked data
func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Requester != nil {
		if err := t.Requester(req); err != nil {
			return nil, err
		}
	}
	rec := new(httptest.ResponseRecorder)
	if t.ContentType != "" {
		rec.Header().Set("Content-Type", t.ContentType)
	}
	if t.Filename != "" {
		rec.Body = new(bytes.Buffer)
		http.ServeFile(rec, req, filepath.Join("testdata", t.Filename))
	}
	res := rec.Result()
	if t.Status > 0 {
		res.StatusCode = t.Status
	}
	if t.Responder != nil {
		if err := t.Responder(res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

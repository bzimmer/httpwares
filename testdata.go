package transport

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Requester .
type Requester func(*http.Request) error

// Responder .
type Responder func(*http.Response) error

// TestDataTransport .
type TestDataTransport struct {
	Status      int
	Filename    string
	ContentType string
	Requester   Requester
	Responder   Responder
}

// RoundTrip .
func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		err  error
		data []byte
	)
	if t.Requester != nil {
		err = t.Requester(req)
		if err != nil {
			return nil, err
		}
	}
	if t.Filename != "" {
		filename := filepath.Join("testdata", t.Filename)
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		data = make([]byte, 0)
	}

	header := make(http.Header)
	header.Set("Content-Type", t.ContentType)

	res := &http.Response{
		StatusCode:    t.Status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}
	if t.Responder != nil {
		err = t.Responder(res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

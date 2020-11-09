package transport

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	// FgRequest is the color of the request
	FgRequest = color.New(color.FgGreen)

	// FgResponse is the color of the response
	FgResponse = color.New(color.FgYellow)
)

// VerboseTransport .
type VerboseTransport struct {
	Writer    io.Writer
	Transport http.RoundTripper
}

func (t *VerboseTransport) isText(header http.Header) bool {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		return false
	}
	if strings.HasPrefix(contentType, "text/") {
		return true
	}
	// content-type is two parts:
	//  - type
	//  - parameters
	splits := strings.Split(contentType, ";")
	switch splits[0] {
	case "application/xml":
	case "application/json":
	case "application/ld+json":
	case "application/geojson":
	default:
		return false
	}
	return true
}

// RoundTrip .
func (t *VerboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	writer := t.Writer
	if writer == nil {
		writer = color.Error
	}
	transport := http.DefaultTransport
	if t.Transport != nil {
		transport = t.Transport
	}
	dump, _ := httputil.DumpRequestOut(req, t.isText(req.Header))
	FgRequest.Fprintln(writer, string(dump))
	res, err := transport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	dump, _ = httputil.DumpResponse(res, t.isText(res.Header))
	FgResponse.Fprintln(writer, string(dump))
	return res, err
}

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

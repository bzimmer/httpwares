package httpwares

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func ContentTypes() []string {
	return []string{
		"application/geo+json",
		"application/geojson",
		"application/gpx+xml",
		"application/json",
		"application/ld+json",
		"application/vnd.google-earth.kml+xml",
		"application/x-www-form-urlencoded",
		"application/xml",
	}
}

// VerboseTransport logs the request and response
type VerboseTransport struct {
	Writer       io.Writer
	Transport    http.RoundTripper
	ContentTypes []string
}

func (t *VerboseTransport) isText(header http.Header) bool {
	var contentTypes = t.ContentTypes
	if contentTypes == nil {
		contentTypes = ContentTypes()
	}
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
	for _, c := range contentTypes {
		if c == splits[0] {
			return true
		}
	}
	return false
}

// RoundTrip is a logging RoundTripper
func (t *VerboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	writer := t.Writer
	if writer == nil {
		writer = os.Stderr
	}
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	out, err := httputil.DumpRequestOut(req, t.isText(req.Header))
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write(out); err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	out, err = httputil.DumpResponse(res, t.isText(res.Header))
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write(out); err != nil {
		return nil, err
	}
	return res, nil
}

package httpwares

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// VerboseTransport logs the request and response
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
	case "application/geo+json":
	case "application/x-www-form-urlencoded":
	default:
		return false
	}
	return true
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
	fmt.Fprintln(writer, string(out))
	res, err := transport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	out, err = httputil.DumpResponse(res, t.isText(res.Header))
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(writer, string(out))
	return res, err
}

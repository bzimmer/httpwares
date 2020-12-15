package httpwares

import (
	"io"
	"net/http"
	"net/http/httputil"
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
	case "application/geo+json":
	case "application/x-www-form-urlencoded":
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

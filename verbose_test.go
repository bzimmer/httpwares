package httpwares_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"
)

func Test_VerboseTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []http.RoundTripper{
		&httpwares.VerboseTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
		&httpwares.VerboseTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
		&httpwares.VerboseTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "transport.json",
				Status:      http.StatusOK,
				ContentType: "application/json; charset=utf-8",
			},
		},
		&httpwares.VerboseTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "",
				Status:      http.StatusNoContent,
				ContentType: "application/foobar; charset=utf-16",
			},
		},
	}
	for i, transport := range tests {
		client := http.Client{
			Transport: transport,
		}
		res, err := client.Get("http://example.com")
		a.NoError(err)
		a.NotNil(res)

		body, err := ioutil.ReadAll(res.Body)
		a.NoError(err)
		a.NotNil(res)

		var quote string
		switch i {
		case 2:
			quote = `{
    "quote": "The mountains are calling & I must go & I will work on while I can, studying incessantly."
}`
		case 3:
			quote = ""
		default:
			quote = "The mountains are calling & I must go & I will work on while I can, studying incessantly."
		}
		a.Equal(quote, strings.Trim(string(body), "\n"))
	}
}

func Test_VerboseLoggingError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client := http.Client{
		Transport: &httpwares.VerboseTransport{
			Transport: &httpwares.TestDataTransport{
				Filename:    "",
				Status:      http.StatusNoContent,
				ContentType: "text/plain",
				Requester: func(req *http.Request) error {
					return errors.New("argh")
				},
			}}}
	res, err := client.Get("http://example.com")
	a.Error(err)
	a.Nil(res)
	a.Equal(`Get "http://example.com": argh`, err.Error())
}

package transport_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/transport"
)

func Test_VerboseTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []http.RoundTripper{
		&transport.VerboseTransport{
			Transport: &transport.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
		&transport.VerboseTransport{
			Transport: &transport.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
		&transport.VerboseTransport{
			Transport: &transport.TestDataTransport{
				Filename:    "transport.json",
				Status:      http.StatusOK,
				ContentType: "application/json; charset=utf-8",
			},
		},
		&transport.VerboseTransport{
			Transport: &transport.TestDataTransport{
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
		Transport: &transport.VerboseTransport{
			Transport: &transport.TestDataTransport{
				Filename:    "",
				Status:      http.StatusNoContent,
				ContentType: "text/plain",
				Requester: func(req *http.Request) error {
					return errors.New("argh")
				},
			}}}
	res, err := client.Get("example.com")
	a.Error(err)
	a.Nil(res)
	a.Equal(`Get "example.com": argh`, err.Error())
}

func Test_TestDataTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := [][]string{
		{"transport.txt", `The mountains are calling & I must go & I will work on while I can, studying incessantly.`},
		{"", ""}}

	for _, test := range tests {
		client := http.Client{
			Transport: &transport.TestDataTransport{
				Filename:    test[0],
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		}
		res, err := client.Get("http://example.com")
		a.NoError(err)
		a.NotNil(res)

		body, err := ioutil.ReadAll(res.Body)
		a.NoError(err)
		a.NotNil(res)
		a.Equal(test[1], strings.Trim(string(body), "\n"))
	}

	client := http.Client{
		Transport: &transport.TestDataTransport{
			Filename:    "~garbage~",
			Status:      http.StatusOK,
			ContentType: "text/plain",
		},
	}
	res, err := client.Get("http://example.com")
	a.Error(err)
	a.Nil(res)
}

func Test_TestRequesterResponder(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := &transport.TestDataTransport{
		Filename:    "transport.txt",
		Status:      http.StatusOK,
		ContentType: "text/plain",
		Requester: func(req *http.Request) error {
			return nil
		},
		Responder: func(res *http.Response) error {
			return nil
		},
	}
	client := http.Client{
		Transport: p,
	}
	res, err := client.Get("http://example.com")
	a.NoError(err)
	a.NotNil(res)

	p.Responder = nil
	p.Requester = func(req *http.Request) error {
		return errors.New("foo")
	}
	res, err = client.Get("http://example.com")
	a.Error(err)
	a.Nil(res)

	p.Requester = nil
	p.Responder = func(res *http.Response) error {
		return errors.New("bar")
	}
	res, err = client.Get("http://example.com")
	a.Error(err)
	a.Nil(res)
}

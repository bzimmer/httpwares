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

func Test_TestDataTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := [][]string{
		{"transport.txt", `The mountains are calling & I must go & I will work on while I can, studying incessantly.`},
		{"", ""}}

	for _, test := range tests {
		client := http.Client{
			Transport: &httpwares.TestDataTransport{
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
		Transport: &httpwares.TestDataTransport{
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

	p := &httpwares.TestDataTransport{
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

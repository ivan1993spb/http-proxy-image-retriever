package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCheckRequest(t *testing.T) {
	request := func(method, path string) *http.Request {
		r, err := http.NewRequest(method, path, nil)
		if err != nil {
			assert.FailNow(t, "request error", err)
		}

		return r
	}

	tests := []*struct {
		r          *http.Request
		statusCode int
	}{
		{request(http.MethodPost, "http://localhost/"), http.StatusMethodNotAllowed},
		{request(http.MethodPut, "http://localhost/"), http.StatusMethodNotAllowed},
		{request(http.MethodTrace, "http://localhost/"), http.StatusMethodNotAllowed},
		{request(http.MethodGet, "http://localhost/1"), http.StatusNotFound},
		{request(http.MethodGet, "http://localhost/2"), http.StatusNotFound},
		{request(http.MethodGet, "http://localhost/3/4/5"), http.StatusNotFound},
		{request(http.MethodGet, "http://localhost/index.html"), http.StatusNotFound},
		{request(http.MethodGet, "http://localhost/"), http.StatusOK},
	}

	var (
		emptyNextHandler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
		emptyLogger      = log.New(ioutil.Discard, "", 0)
		handler          = MiddlewareCheckRequest(emptyNextHandler, emptyLogger)
	)

	for _, test := range tests {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, test.r)
		assert.Equal(t, test.statusCode, w.Code, "status code")
	}
}

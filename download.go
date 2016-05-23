package main

import (
	"flag"
	"net/http"
	"net/url"
)

var HTTPClient = &http.Client{}

func init() {
	flag.DurationVar(&HTTPClient.Timeout, "timeout", 0, "time limit for requests")
}

func Download(URL *url.URL, stopChan <-chan struct{}) (<-chan *http.Response, <-chan error) {
	var (
		respChan = make(chan *http.Response)
		errChan  = make(chan error)
	)

	go func() {
		resp, err := HTTPClient.Do(&http.Request{
			URL:    URL,
			Cancel: stopChan,
		})

		if err != nil {
			errChan <- err
		} else {
			respChan <- resp
		}

		close(respChan)
		close(errChan)
	}()

	return respChan, errChan
}

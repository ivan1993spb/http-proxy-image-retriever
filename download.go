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

func download(URL string, stopChan <-chan struct{}) (<-chan *http.Response, <-chan error) {
	var (
		respChan = make(chan *http.Response)
		errChan  = make(chan error)
	)

	go func() {
		if parsedURL, err := url.Parse(URL); err != nil {
			errChan <- err
		} else {
			resp, err := HTTPClient.Do(&http.Request{
				URL:    parsedURL,
				Cancel: stopChan,
			})
			if err != nil {
				errChan <- err
			} else {
				respChan <- resp
			}
		}

		close(respChan)
		close(errChan)
	}()

	return respChan, errChan
}

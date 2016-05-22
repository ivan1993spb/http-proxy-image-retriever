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

func download(URL *url.URL, stopChan <-chan struct{}) (*http.Response, error) {
	return HTTPClient.Do(&http.Request{
		URL:    URL,
		Cancel: stopChan,
	})
}

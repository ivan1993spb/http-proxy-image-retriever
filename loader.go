package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

type Loader struct {
	logger *log.Logger
	queue  chan struct{}
	*http.Client
}

func NewLoader(logger *log.Logger, maxWorkerCount uint, timeout time.Duration) *Loader {
	return &Loader{
		logger: logger,
		queue:  make(chan struct{}, maxWorkerCount),
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (l *Loader) DownloadCallback(stopChan <-chan struct{}, URL *url.URL, callback func(*http.Response, error)) {
	go func() {
		l.logger.Println("loading url:", URL)

		select {
		case l.queue <- struct{}{}:

			callback(l.Client.Do(&http.Request{
				URL:    URL,
				Cancel: stopChan,
			}))

			<-l.queue
		case <-stopChan:
		}

		l.logger.Println("loader finished:", URL)
	}()
}

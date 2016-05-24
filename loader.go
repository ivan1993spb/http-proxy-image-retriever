package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	RequestTimeout time.Duration
	MaxWorkerCount uint
)

func init() {
	flag.DurationVar(&RequestTimeout, "timeout", 0, "time limit for requests")
	flag.UintVar(&MaxWorkerCount, "max-worker-count", 500, "max count of workers")
}

type Loader struct {
	logger *log.Logger
	queue  chan struct{}
	*http.Client
}

func NewLoader(logger *log.Logger) *Loader {
	return &Loader{
		logger: logger,
		queue:  make(chan struct{}, MaxWorkerCount),
		Client: &http.Client{
			Timeout: RequestTimeout,
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
				Close:  true,
				Cancel: stopChan,
			}))

			<-l.queue
		case <-stopChan:
		}

		l.logger.Println("loading finished:", URL)
	}()
}

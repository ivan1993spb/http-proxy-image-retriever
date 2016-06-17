package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	requestTimeout time.Duration
	maxWorkerCount uint
)

func init() {
	flag.DurationVar(&requestTimeout, "timeout", time.Second*15, "time limit for requests")
	flag.UintVar(&maxWorkerCount, "max-worker-count", 500, "max count of workers")
}

// Loader is http client with queue, timeouts and cancel chan
type Loader struct {
	logger *log.Logger
	queue  chan struct{}
	*http.Client
}

// NewLoader creates new Loader
func NewLoader(logger *log.Logger, cancel <-chan struct{}) *Loader {
	return &Loader{
		logger: logger,
		queue:  make(chan struct{}, maxWorkerCount),
		Client: &http.Client{
			Timeout: requestTimeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout: 30 * time.Second,
					Cancel:  cancel,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true,
			},
		},
	}
}

// DownloadCallback runs request concurrently and calls passed callback func for response handling
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

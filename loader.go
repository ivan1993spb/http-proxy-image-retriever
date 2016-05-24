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

func NewLoader(logger *log.Logger, cancel <-chan struct{}) *Loader {
	return &Loader{
		logger: logger,
		queue:  make(chan struct{}, MaxWorkerCount),
		Client: &http.Client{
			Timeout: RequestTimeout,
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

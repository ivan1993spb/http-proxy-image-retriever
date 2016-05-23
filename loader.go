package main

import (
	"log"
	"net/http"
	"net/url"
	//"sync"
)

type Loader struct {
	logger *log.Logger
	queue  chan struct{}
	*http.Client
}

func NewLoader() *Loader {
	return nil
}

func (l *Loader) Download(stopChan <-chan struct{}, URL *url.URL, callback func(*http.Response, error)) {
	//l.waitGroup.Add(1)

	go func() {
		select {
		case l.queue <- struct{}{}:

			callback(l.Client.Do(&http.Request{
				URL:    URL,
				Cancel: stopChan,
			}))

			<-l.queue
		case <-stopChan:
		}

		//l.waitGroup.Done()
	}()

}

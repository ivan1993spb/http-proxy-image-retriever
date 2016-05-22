package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type HTTPImageProxyHandler struct {
	logger   *log.Logger
	stopLock sync.Mutex
	stopChan chan struct{}
}

func NewHTTPImageProxyHandler(logger *log.Logger) *HTTPImageProxyHandler {
	return &HTTPImageProxyHandler{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (h *HTTPImageProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	stopChan := h.getRequestStopChan(w)

	h.logger.Println("processing url:", r.FormValue("url"))

	loadedChan := make(chan struct{})
	go func() {
		URL, _ := url.Parse(r.FormValue("url"))
		download(URL, stopChan)
		close(loadedChan)
		fmt.Println("stop loading")
	}()

	select {
	case <-stopChan:
		fmt.Println("interrupted")
	case <-loadedChan:
		fmt.Println("loaded")
	}

}

// getRequestStopChan returns stop chan for a request
func (h *HTTPImageProxyHandler) getRequestStopChan(w http.ResponseWriter) <-chan struct{} {
	if cn, ok := w.(http.CloseNotifier); ok {
		stopChan := make(chan struct{})

		go func() {
			select {
			case <-h.stopChan:
			case <-cn.CloseNotify():
			}
			close(stopChan)
		}()

		return stopChan
	}

	return h.stopChan
}

func (h *HTTPImageProxyHandler) Stop() {
	h.stopLock.Lock()
	defer h.stopLock.Unlock()

	if h.stopChan != nil {
		close(h.stopChan)
		h.stopChan = nil
	}
}

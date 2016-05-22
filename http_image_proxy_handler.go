package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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

	select {
	case <-h.getRequestStopChan(w):
		fmt.Println("okok")
	case <-time.After(time.Second * 10):
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

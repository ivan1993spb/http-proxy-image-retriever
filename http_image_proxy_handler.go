package main

import (
	"fmt"
	"log"
	"net/http"
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
	h.logger.Println("processing url:", r.FormValue("url"))

	stopChan := h.getRequestStopChan(w)
	respChan, errChan := download(r.FormValue("url"), stopChan)

	select {
	case <-stopChan:
		// Connection closed or processing interrupted
	case err := <-errChan:
		h.logger.Println("loading error:", err)
		HTTPErrorHTML(w, "cannot load url", http.StatusOK)
	case resp := <-respChan:
		fmt.Println("received response")
		resp.Body
	}
}

// getRequestStopChan returns stop chan for a request
func (h *HTTPImageProxyHandler) getRequestStopChan(w http.ResponseWriter) <-chan struct{} {
	if cn, ok := w.(http.CloseNotifier); ok {
		stopChan := make(chan struct{})

		go func() {
			select {
			case <-h.stopChan:
				h.logger.Println("processing interrupted")
			case <-cn.CloseNotify():
				h.logger.Println("connection closed")
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

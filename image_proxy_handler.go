package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
)

// ImageProxyHandler accepts http request with url param, downloads
// html page from passed url, parses html and finds all images, downloads
// all found images, generates response html page with found images included
// into page by data URI scheme.
type ImageProxyHandler struct {
	logger   *log.Logger
	stopLock sync.Mutex
	stopChan chan struct{}
}

// NewImageProxyHandler creates new ImageProxyHandler
func NewImageProxyHandler(logger *log.Logger) *ImageProxyHandler {
	return &ImageProxyHandler{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Implementing http.Handler interface
func (h *ImageProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	URL, err := url.Parse(r.FormValue("url"))
	if err != nil {
		h.logger.Println("invalid url:", err)
		HTTPErrorHTML(w, "invalid url", http.StatusOK)
		return
	}

	h.logger.Println("processing url:", URL)

	stopChan := h.getRequestStopChan(w)
	respChan, errChan := Download(URL, stopChan)

	select {
	case <-stopChan:
		// Connection closed or processing interrupted
	case err := <-errChan:
		h.logger.Println("loading error:", err)
		HTTPErrorHTML(w, "cannot load url", http.StatusOK)
		return
	case resp := <-respChan:
		h.logger.Println("received response")
		imgSources, err := FindImageSources(stopChan, resp.Body)
		if err != nil {
			h.logger.Println("parsing response error:", err)
			HTTPErrorHTML(w, "cannot parse loaded html page", http.StatusOK)
			return
		}

		for _, source := range imgSources {
			// TODO fix image url: `path/to/image.png`, `/path/to/image.png`, `http://ex.ple/path/to/image.png`
			// TODO add URL data case
			h.logger.Println("found src:", source)
		}
	}
}

// getRequestStopChan returns stop chan for a request
func (h *ImageProxyHandler) getRequestStopChan(w http.ResponseWriter) <-chan struct{} {
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

// Stop stops all processing goroutines started by handler
func (h *ImageProxyHandler) Stop() {
	h.stopLock.Lock()
	defer h.stopLock.Unlock()

	if h.stopChan != nil {
		close(h.stopChan)
		h.stopChan = nil
	}
}

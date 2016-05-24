package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/vincent-petithory/dataurl"
)

// ImageProxyHandler accepts http request with url param, downloads
// html page from passed url, parses html and finds all images, downloads
// all found images, generates response html page with found images included
// into page by data URI scheme.
type ImageProxyHandler struct {
	logger   *log.Logger
	loader   *Loader
	stopLock sync.Mutex
	stopChan chan struct{}
}

// NewImageProxyHandler creates new ImageProxyHandler
func NewImageProxyHandler(logger *log.Logger) *ImageProxyHandler {
	return &ImageProxyHandler{
		logger:   logger,
		loader:   NewLoader(logger, WorkerCount, RequestTimeout),
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
	imageURLChan, dataURLChan1, errorChan1 := h.findImagesPageURL(stopChan, URL)
	dataURLChan2, errorChan2 := h.loadImages(stopChan, imageURLChan)

	errorChan := MergeErrorChans(stopChan, errorChan1, errorChan2)
	dataURLChan := MergeDataURLChans(stopChan, dataURLChan1, dataURLChan2)

	go func() {
		for {
			select {
			case err := <-errorChan:
				h.logger.Println("error", err)
			case <-stopChan:
				return
			}
		}
	}()

	ImagesHTMLRender(w, dataURLChan)
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

type ErrFindImagesPageURL struct {
	err string
}

func (e *ErrFindImagesPageURL) Error() string {
	return "find images on page error: " + e.err
}

func (h *ImageProxyHandler) findImagesPageURL(stopChan <-chan struct{}, URL *url.URL) (
	<-chan *url.URL, <-chan *dataurl.DataURL, <-chan error) {

	var (
		imageURLChan = make(chan *url.URL)
		dataURLChan  = make(chan *dataurl.DataURL)
		errorChan    = make(chan error)
	)

	h.loader.DownloadCallback(stopChan, URL, func(resp *http.Response, err error) {
		defer func() {
			close(imageURLChan)
			close(dataURLChan)
			close(errorChan)
		}()

		if err != nil {
			errorChan <- &ErrFindImagesPageURL{"loading html page error: " + err.Error()}
			return
		}

		imgSources, err := FindImageSources(resp.Body)
		if err != nil {
			errorChan <- &ErrFindImagesPageURL{"parsing response error: " + err.Error()}
			return
		}

		for _, source := range imgSources {
			if IsDataUrl(source) {
				h.logger.Println("found dataurl image source")
				if du, err := dataurl.DecodeString(source); err != nil {
					errorChan <- &ErrFindImagesPageURL{"cannot decode dataurl: " + err.Error()}
				} else {
					dataURLChan <- du
				}
			} else {
				h.logger.Println("found link image source")
				if imageURL, err := url.Parse(source); err != nil {
					errorChan <- &ErrFindImagesPageURL{"cannot parse image url: " + err.Error()}
				} else {
					imageURLChan <- resp.Request.URL.ResolveReference(imageURL)
				}
			}
		}
	})

	return imageURLChan, dataURLChan, errorChan
}

type ErrLoadingImage struct {
	err string
}

func (e *ErrLoadingImage) Error() string {
	return "loading image error: " + e.err
}

func (h *ImageProxyHandler) loadImages(stopChan <-chan struct{}, imageURLChan <-chan *url.URL) (
	<-chan *dataurl.DataURL, <-chan error) {

	var (
		dataURLChan = make(chan *dataurl.DataURL)
		errorChan   = make(chan error)
	)

	go func() {
		var wg sync.WaitGroup

		for URL := range imageURLChan {
			wg.Add(1)
			h.loader.DownloadCallback(stopChan, URL, func(resp *http.Response, err error) {
				if err != nil {
					errorChan <- &ErrLoadingImage{err.Error()}
				} else if resp.StatusCode != http.StatusOK {
					errorChan <- &ErrLoadingImage{"bad status code"}
				} else if contentType := resp.Header.Get("Content-Type"); !IsBrowserImageContentType(contentType) {
					errorChan <- &ErrLoadingImage{"unexpected content-type: " + contentType}
				} else if data, err := ioutil.ReadAll(resp.Body); err != nil {
					errorChan <- &ErrLoadingImage{err.Error()}
				} else {
					h.logger.Println("image loaded", contentType)
					dataURLChan <- dataurl.New(data, contentType)
				}

				wg.Done()
			})
		}

		go func() {
			wg.Wait()
			h.logger.Println("closing dataURL chan")
			close(dataURLChan)
			close(errorChan)
		}()
	}()

	return dataURLChan, errorChan
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

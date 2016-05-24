package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/vincent-petithory/dataurl"
)

// Buffer sizes
const (
	IMAGE_URL_CHAN_BUFFER_SIZE = 20
	DATAURL_CHAN_BUFFER_SIZE   = 20
	ERROR_CHAN_BUFFER_SIZE     = 10
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
	stopChan := make(chan struct{})

	return &ImageProxyHandler{
		logger:   logger,
		loader:   NewLoader(logger, stopChan),
		stopChan: stopChan,
	}
}

// Implementing http.Handler interface
func (h *ImageProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	URL, err := url.Parse(r.FormValue("url"))
	if err != nil {
		h.logger.Println("invalid url:", err)
		ErrorHTML(w, "invalid url", http.StatusOK)
		return
	}

	h.logger.Println("processing url:", URL)

	stopChan := h.getRequestStopChan(w)
	imageURLChan, dataURLChan1, errorChan1 := h.findImagesPageURL(stopChan, URL)
	dataURLChan2, errorChan2 := h.loadImages(stopChan, imageURLChan)
	errorChan := MergeErrorChans(stopChan, errorChan1, errorChan2)
	dataURLChan := MergeDataURLChans(stopChan, dataURLChan1, dataURLChan2)

	go func() {
		for err := range errorChan {
			h.logger.Println("request handling error:", err)
		}
	}()

	h.imagesHTML(w, dataURLChan)
}

// getRequestStopChan returns stop chan for a request
func (h *ImageProxyHandler) getRequestStopChan(w http.ResponseWriter) <-chan struct{} {
	if cn, ok := w.(http.CloseNotifier); ok {
		stopChan := make(chan struct{})
		closeConnChan := cn.CloseNotify()

		go func() {
			select {
			case <-h.stopChan:
				h.logger.Println("processing interrupted")
			case <-closeConnChan:
				h.logger.Println("connection closed")
			}
			close(stopChan)
		}()

		return stopChan
	}

	return h.stopChan
}

type ErrFindImagesPageURL string

func (e ErrFindImagesPageURL) Error() string {
	return "finding images on page error: " + string(e)
}

func (h *ImageProxyHandler) findImagesPageURL(stopChan <-chan struct{}, URL *url.URL) (
	<-chan *url.URL, <-chan *dataurl.DataURL, <-chan error) {

	var (
		imageURLChan = make(chan *url.URL, IMAGE_URL_CHAN_BUFFER_SIZE)
		dataURLChan  = make(chan *dataurl.DataURL, DATAURL_CHAN_BUFFER_SIZE)
		errorChan    = make(chan error, ERROR_CHAN_BUFFER_SIZE)
	)

	h.loader.DownloadCallback(stopChan, URL, func(resp *http.Response, err error) {
		defer func() {
			close(imageURLChan)
			close(dataURLChan)
			close(errorChan)
		}()

		if err != nil {
			errorChan <- ErrFindImagesPageURL("loading html page error: " + err.Error())
			return
		}

		imgSources, err := FindImageSources(resp.Body)
		if err != nil {
			errorChan <- ErrFindImagesPageURL("parsing response error: " + err.Error())
			return
		}

		h.logger.Println("image sources found:", len(imgSources))

		for _, source := range imgSources {
			if IsDataUrl(source) {
				h.logger.Println("found dataurl image source")
				if du, err := dataurl.DecodeString(source); err != nil {
					errorChan <- ErrFindImagesPageURL("cannot decode dataurl: " + err.Error())
				} else {
					dataURLChan <- du
				}
			} else {
				h.logger.Println("found link image source")
				if imageURL, err := url.Parse(source); err != nil {
					errorChan <- ErrFindImagesPageURL("cannot parse image url: " + err.Error())
				} else {
					imageURLChan <- resp.Request.URL.ResolveReference(imageURL)
				}
			}
		}
	})

	return imageURLChan, dataURLChan, errorChan
}

type ErrLoadingImage string

func (e ErrLoadingImage) Error() string {
	return "loading image error: " + string(e)
}

func (h *ImageProxyHandler) loadImages(stopChan <-chan struct{}, imageURLChan <-chan *url.URL) (
	<-chan *dataurl.DataURL, <-chan error) {

	var (
		dataURLChan = make(chan *dataurl.DataURL, DATAURL_CHAN_BUFFER_SIZE)
		errorChan   = make(chan error, ERROR_CHAN_BUFFER_SIZE)
	)

	go func() {
		var wg sync.WaitGroup

		for URL := range imageURLChan {
			wg.Add(1)
			h.loader.DownloadCallback(stopChan, URL, func(resp *http.Response, err error) {
				if err != nil {
					errorChan <- ErrLoadingImage(err.Error())
				} else if resp.StatusCode != http.StatusOK {
					errorChan <- ErrLoadingImage("bad status code")
				} else if contentType := resp.Header.Get("Content-Type"); !IsBrowserImageMIME(contentType) {
					errorChan <- ErrLoadingImage("unexpected content-type: " + contentType)
				} else if data, err := ioutil.ReadAll(resp.Body); err != nil {
					errorChan <- ErrLoadingImage(err.Error())
				} else {
					h.logger.Println("image loaded", contentType)
					dataURLChan <- dataurl.New(data, contentType)
				}
				wg.Done()
			})
		}

		go func() {
			wg.Wait()
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

var ImagesPageTmpl = template.Must(template.New("images_page").Parse(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Images</title>
        <style>
            img { display: block; margin: 10px; }
        </style>
    </head>
    <body>
        <h1>Images</h1>
        {{range .}}<img src="{{html .}}">
        {{else}}<b>No images</b>{{end}}
    </body>
</html>
`))

// imagesHTML sends html with found images
func (h *ImageProxyHandler) imagesHTML(w http.ResponseWriter, dataURLChan <-chan *dataurl.DataURL) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := ImagesPageTmpl.Execute(w, dataURLChan); err != nil {
		h.logger.Println("writing html response error:", err)
	}
}

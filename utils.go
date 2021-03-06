package main

import (
	"html/template"
	"net/http"
	"regexp"
	"sync"

	"github.com/vincent-petithory/dataurl"
)

var expDetectDataURL = regexp.MustCompile(
	`(?i)^\s*data:([a-z]+\/[a-z0-9\-\+]+(;[a-z\-]+\=[a-z0-9\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$`)

// IsDataUrl returns true if passed string s is data url string
func IsDataUrl(s string) bool {
	return expDetectDataURL.MatchString(s)
}

// IsBrowserImageMIME returns true if passed string mime is valid mime
// type for images in web browsers
func IsBrowserImageMIME(mime string) bool {
	switch mime {
	case "image/jpeg", "image/jp2", "image/jpx", "image/jpm", "image/webp", "image/vnd.ms-photo",
		"image/jxr", "image/gif", "image/png", "image/tiff", "image/tiff-fx", "image/svg+xml",
		"image/x‑xbitmap", "image/x‑xbm", "image/bmp", "image/x-bmp", "image/x-icon":
		return true
	}

	return false
}

var errorPageTmpl = template.Must(template.New("error_page").Parse(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Error</title>
    </head>
    <body>
        <h1>Error</h1>
        <p>{{.}}</p>
    </body>
</html>
`))

// ErrorHTML sends error message with specific status code
func ErrorHTML(w http.ResponseWriter, error string, code int) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	return errorPageTmpl.Execute(w, error)
}

func mergeErrorChans(stopChan <-chan struct{}, errorChans ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	errorChanOut := make(chan error, ERROR_CHAN_BUFFER_SIZE)

	output := func(errorChan <-chan error) {
		for err := range errorChan {
			select {
			case errorChanOut <- err:
			case <-stopChan:
			}
		}
		wg.Done()
	}

	wg.Add(len(errorChans))
	for _, errorChan := range errorChans {
		go output(errorChan)
	}

	go func() {
		wg.Wait()
		close(errorChanOut)
	}()

	return errorChanOut
}

func mergeDataURLChans(stopChan <-chan struct{}, dataURLChans ...<-chan *dataurl.DataURL) <-chan *dataurl.DataURL {
	var wg sync.WaitGroup
	dataURLChanOut := make(chan *dataurl.DataURL, DATAURL_CHAN_BUFFER_SIZE)

	output := func(dataURLChan <-chan *dataurl.DataURL) {
		for dataURL := range dataURLChan {
			select {
			case dataURLChanOut <- dataURL:
			case <-stopChan:
			}
		}
		wg.Done()
	}

	wg.Add(len(dataURLChans))
	for _, dataURLChan := range dataURLChans {
		go output(dataURLChan)
	}

	go func() {
		wg.Wait()
		close(dataURLChanOut)
	}()

	return dataURLChanOut
}

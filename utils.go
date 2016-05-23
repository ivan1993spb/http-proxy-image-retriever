package main

import (
	"github.com/vincent-petithory/dataurl"
	"html/template"
	"net/http"
	"regexp"
)

var ExpDetectDataURL = regexp.MustCompile(
	`(?i)^\s*data:([a-z]+\/[a-z0-9\-\+]+(;[a-z\-]+\=[a-z0-9\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$`)

// IsDataUrl returns true if passed string s is data url
func IsDataUrl(s string) bool {
	return ExpDetectDataURL.MatchString(s)
}

func IsBrowserImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jp2", "image/jpx", "image/jpm", "image/webp", "image/vnd.ms-photo",
		"image/jxr", "image/gif", "image/png", "image/tiff", "image/tiff-fx", "image/svg+xml",
		"image/x‑xbitmap", "image/x‑xbm", "image/bmp", "image/x-bmp", "image/x-icon":
		return true
	}

	return false
}

var ErrorPageTmpl = template.Must(template.New("error_page").Parse(`<!DOCTYPE html>
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

// HTTPErrorHTML sends error message with specific status code
func HTTPErrorHTML(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	ErrorPageTmpl.Execute(w, error)
}

var ImagesPageTmpl = template.Must(template.New("images_page").Parse(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Images</title>
    </head>
    <body>
        <h1>Images</h1>
        {{range .}}<img src="{{.}}">
        {{else}}<b>No images</b>{{end}}
    </body>
</html>
`))

func ImagesHTMLRender(w http.ResponseWriter, dataURLChan <-chan *dataurl.DataURL) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ImagesPageTmpl.Execute(w, dataURLChan)
}

func MergeErrorChans(stopChan <-chan struct{}, errorChan1 <-chan error, errorChan2 <-chan error) <-chan error {
	errorChan := make(chan error)

	go func() {
		defer close(errorChan)

		select {
		case <-stopChan:
			return
		case errorChan <- <-errorChan1:
		case errorChan <- <-errorChan2:
		}
	}()

	return errorChan
}

func MergeDataURLChans(stopChan <-chan struct{}, dataURLChan1 <-chan *dataurl.DataURL, dataURLChan2 <-chan *dataurl.DataURL) <-chan *dataurl.DataURL {
	dataURLChan := make(chan *dataurl.DataURL)

	go func() {
		defer close(dataURLChan)

		select {
		case <-stopChan:
			return
		case dataURLChan <- <-dataURLChan1:
		case dataURLChan <- <-dataURLChan2:
		}
	}()

	return dataURLChan
}

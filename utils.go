package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"regexp"
	//"strings"
	"sync"

	"github.com/vincent-petithory/dataurl"
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
        {{range .}}<img src="{{html .}}">
        {{else}}<b>No images</b>{{end}}
    </body>
</html>
`))

func ImagesHTMLRender(w http.ResponseWriter, dataURLChan <-chan *dataurl.DataURL) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//dataURLs := []string{}
	//for dataURL := range dataURLChan {
	//	dataURL.String()
	//	dataURLs = append(dataURLs, dataURL.String())
	//	//dataURLs = append(dataURLs, strings.Repeat("A", 111500))
	//}
	//fmt.Println(dataURLs)
	ImagesPageTmpl.Execute(w, dataURLChan)
}

func MergeErrorChans(stopChan <-chan struct{}, errorChans ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	errorChanOut := make(chan error)

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

func MergeDataURLChans(stopChan <-chan struct{}, dataURLChans ...<-chan *dataurl.DataURL) <-chan *dataurl.DataURL {
	var wg sync.WaitGroup
	dataURLChanOut := make(chan *dataurl.DataURL)

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

package main

import (
	"log"
	"net/http"
	"net/url"
)

type Image struct {
	MimeType   string
	dataBase64 []byte
}

type ErrLoading struct {
	err string
}

func (err *ErrLoading) Error() string {
	return "loading error: " + err.err
}

func DownloadImages(stopChan <-chan struct{}, urlChan <-chan *url.URL) (<-chan *Image, <-chan error) {
	imageChan := make(chan *Image)
	errChan := make(chan error)

	go func() {
		defer close(imageChan)
		defer close(errChan)

		for {
			select {
			case URL := <-urlChan:
				respChan, downloadErrChan := Download(URL, stopChan)
				select {
				case <-stopChan:
					return
				case resp := <-respChan:
					if resp.StatusCode == http.StatusOK {
						errChan <- &ErrLoading{"unexpected status: " + resp.Status}
					} else if !IsBrowserImageContentType(resp.Header.Get("Content-Type")) {
						errChan <- &ErrLoading{"unexpected content-type"}
					} else {

					}
				case errChan <- &ErrLoading{(<-downloadErrChan).Error()}:
				}
			case <-stopChan:
				return
			}
		}
	}()

	return imageChan, errChan
}

func LogErrorChan(stopChan <-chan struct{}, errChanChan <-chan <-chan error, logger *log.Logger) {
	var listenErrChan = func(errChan <-chan error) {
		for err := range errChan {
			logger.Println(err)
			select {
			case <-stopChan:
				return
			default:
			}
		}
	}

	go func() {
		for {
			select {
			case <-stopChan:
				return
			case err := <-errChanChan:
				go listenErrChan(err)
			}
		}
	}()
}

package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Image struct {
	MimeType   string
	Base64Data []byte
}

type ErrLoading struct {
	err string
}

func (err *ErrLoading) Error() string {
	return "loading error: " + err.err
}

func DownloadImages(stopChan <-chan struct{}, urlChan <-chan *url.URL) (<-chan *Image, <-chan error) {
	// TODO getting images from sync pool

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
					} else if mime := resp.Header.Get("Content-Type"); !IsBrowserImageContentType(mime) {
						errChan <- &ErrLoading{"unexpected content-type"}
					} else if rawData, err := ioutil.ReadAll(resp.Body); err != nil {
						errChan <- &ErrLoading{err.Error()}
					} else {
						base64Data := make([]byte, base64.URLEncoding.EncodedLen(len(rawData)))
						base64.URLEncoding.Encode(base64Data, rawData)
						imageChan <- &Image{mime, base64Data}
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

func MergeImageChans(stopChan <-chan struct{}, imageChanChan <-chan <-chan *Image) <-chan *Image {
	outImageChan := make(chan *Image)

	go func() {
		var wg sync.WaitGroup

		output := func(imageChan <-chan *Image) {
			for image := range imageChan {
				select {
				case outImageChan <- image:
				case <-stopChan:
				}
			}
			wg.Done()
		}

		for imageChan := range imageChanChan {
			wg.Add(1)
			go output(imageChan)
		}

		wg.Wait()
		close(outImageChan)
	}()

	return outImageChan
}

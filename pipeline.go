package main

import (
	"io"
	"net/url"
)

func DownloadImages(stopChan <-chan struct{}, urlChan <-chan url.URL) <-chan io.Reader {
	imageChan := make(chan io.Reader)

	go func() {
		close(imageChan)
	}()

	return imageChan
}

func mergeImageChans(stopChan <-chan struct{}, imageChan ...<-chan io.Reader) {

}

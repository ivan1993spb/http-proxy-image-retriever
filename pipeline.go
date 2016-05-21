package main

import (
	"io"
	"net/url"
)

func retrieveUrls(stopChan <-chan struct{}, html io.Reader) <-chan url.URL {
	urlChan := make(chan url.URL)

	go func() {
		// TODO fix image url
		close(urlChan)
	}()

	return urlChan
}

func downloadImages(stopChan <-chan struct{}, urlChan <-chan url.URL) <-chan io.Reader {
	imageChan := make(chan io.Reader)

	go func() {
		close(imageChan)
	}()

	return imageChan
}

func mergeImageChans(stopChan <-chan struct{}, imageChan ...<-chan io.Reader) {

}

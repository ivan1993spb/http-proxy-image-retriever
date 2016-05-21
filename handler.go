package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPProxyHandler struct {
	logger   *log.Logger
	stopChan <-chan struct{}
}

func (h *HTTPProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("accepted request")

	select {
	case <-h.stopChan:
		fmt.Println("okok")
	case <-time.After(time.Second * 10):
	}

	log.Println("finished connection handling")
}

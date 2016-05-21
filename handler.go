package main

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPProxyHandler struct {
}

func (h *HTTPProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accepted request")

	cn, ok := w.(http.CloseNotifier)
	if ok {
		select {
		case <-cn.CloseNotify():
			fmt.Println("okok")
		case <-time.After(time.Second * 10):
		}
	} else {
		<-time.After(time.Second * 10)
	}

	fmt.Println("finished handling")
}

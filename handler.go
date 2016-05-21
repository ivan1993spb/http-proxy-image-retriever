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
	time.Sleep(time.Second * 10)
	fmt.Println("finished handling")

}

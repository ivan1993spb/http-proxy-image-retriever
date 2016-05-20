package main

import (
	"net/http"
)

type HTTPProxyHandler struct {
}

func (h *HTTPProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

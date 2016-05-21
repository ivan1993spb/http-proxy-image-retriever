package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//go:generate go-bindata static/

type HTTPImageProxyHandler struct {
	logger   *log.Logger
	stopChan <-chan struct{}
}

func (h *HTTPImageProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("accepted request")

	if r.Method == http.MethodGet {
		if r.URL.Path == "/" {
			if url := r.FormValue("url"); url != "" {
				cn, ok := w.(http.CloseNotifier)
				var stopChan chan struct{}
				if ok {
					stopChan = make(chan struct{})
					go func() {
						select {
						case <-h.stopChan:
							fmt.Println(1)
						case <-cn.CloseNotify():
							fmt.Println(2)
						}
						close(stopChan)
					}()
				} else {
					go func() {
						<-h.stopChan
						close(stopChan)
					}()
				}

				select {
				case <-stopChan:
					fmt.Println("okok")
				case <-time.After(time.Second * 10):
				}
			} else {
				//data := MustAsset("static/error.html")
				h.logger.Println("passed empty url")
			}
		} else {
			h.logger.Println("invalid path")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	} else {
		h.logger.Println("request method not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	log.Println("finished connection handling")
}

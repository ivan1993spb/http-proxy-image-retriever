package main

import (
	"log"
	"net/http"
	"sync"
)

type HTTPImageProxyHandler struct {
	logger   *log.Logger
	stopLock sync.Mutex
	stopChan chan struct{}
}

func NewHTTPImageProxyHandler(logger *log.Logger) *HTTPImageProxyHandler {
	return &HTTPImageProxyHandler{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (h *HTTPImageProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("accepted request")
	defer log.Println("finished connection handling")

	//cn, ok := w.(http.CloseNotifier)
	//var stopChan chan struct{}
	//if ok {
	//	stopChan = make(chan struct{})
	//	go func() {
	//		select {
	//		case <-h.stopChan:
	//			fmt.Println(1)
	//		case <-cn.CloseNotify():
	//			fmt.Println(2)
	//		}
	//		close(stopChan)
	//	}()
	//} else {
	//	go func() {
	//		<-h.stopChan
	//		close(stopChan)
	//	}()
	//}
	//
	//select {
	//case <-stopChan:
	//	fmt.Println("okok")
	//case <-time.After(time.Second * 10):
	//}

}

func (h *HTTPImageProxyHandler) Stop() {
	h.stopLock.Lock()
	defer h.stopLock.Unlock()

	if h.stopChan != nil {
		close(h.stopChan)
		h.stopChan = nil
	}
}

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	stoppableListener, err := NewStoppableListener(listener)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening on", stoppableListener.Addr())

	go func() {
		err := http.Serve(stoppableListener, &HTTPProxyHandler{})
		log.Println(err)
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println("ok", <-ch)

	stoppableListener.Stop()

	time.Sleep(time.Second)
}

package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Listen on 127.0.0.1:48879.  That's my favorite port number because in
	// hex 48879 is 0xBEEF.
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:48879")
	if nil != err {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())

	// Make a new service and send it into the background.
	service := NewService()
	go service.Serve(listener)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	service.Stop()

}

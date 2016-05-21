package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

var serverAddr string

func init() {
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8888", "server address")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("initializing server")

	stopHandlingChan := make(chan struct{})

	server := &graceful.Server{
		Timeout: time.Second,
		Server: &http.Server{
			Addr: serverAddr,
			Handler: &HTTPImageProxyHandler{
				logger,
				stopHandlingChan,
			},
		},
		Logger: logger,
		ShutdownInitiated: func() {
			close(stopHandlingChan)
		},
	}

	logger.Println("starting server")

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger.Fatal(err)
		}
	}

	logger.Println("server stopped")
}

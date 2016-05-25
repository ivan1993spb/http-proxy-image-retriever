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

var ServerAddr string

func init() {
	flag.StringVar(&ServerAddr, "addr", ":8888", "server address")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("initializing server")

	handler := NewImageProxyHandler(logger)

	server := &graceful.Server{
		Timeout: time.Second,
		Server: &http.Server{
			Addr:    ServerAddr,
			Handler: MiddlewareCheckRequest(handler, logger),
		},
		Logger: logger,
		ShutdownInitiated: func() {
			logger.Println("stop all processing goroutines")
			handler.Stop()
		},
	}

	logger.Println("starting server")

	if err := server.ListenAndServe(); err != nil {
		// If err is critical error
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger.Fatal(err)
		}
	}

	logger.Println("server stopped")
}

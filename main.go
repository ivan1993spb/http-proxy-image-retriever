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

var (
	serverAddr  string
	killTimeout time.Duration
)

func init() {
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8888", "server address")
	flag.DurationVar(&killTimeout, "kill-timeout", time.Second,
		"the duration to allow outstanding requests to survive before forcefully terminating them")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	server := &graceful.Server{
		Timeout: killTimeout,
		Server: &http.Server{
			Addr:    serverAddr,
			Handler: &HTTPProxyHandler{logger},
		},
		Logger: logger,
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger.Fatal(err)
		}
	}

	time.Sleep(time.Second)
}

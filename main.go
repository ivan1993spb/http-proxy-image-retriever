package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

var (
	ServerAddr     string
	RequestTimeout time.Duration
	WorkerCount    uint
)

func init() {
	flag.StringVar(&ServerAddr, "addr", "127.0.0.1:8888", "server address")
	flag.DurationVar(&RequestTimeout, "timeout", 0, "time limit for requests")
	flag.UintVar(&WorkerCount, "worker-count", 500, "count of workers")
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

	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Println("starting server")

	if err := server.ListenAndServe(); err != nil {
		// If err is critical error
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger.Fatal(err)
		}
	}

	logger.Println("server stopped")
}

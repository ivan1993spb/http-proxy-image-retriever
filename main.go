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

var ServerAddr string

func init() {
	flag.StringVar(&ServerAddr, "addr", "127.0.0.1:8888", "server address")
}

func main() {
	flag.Parse()

	//log.Println(net.LookupAddr(ServerAddr))
	//log.Println(net.LookupHost(ServerAddr))
	//log.Println(net.LookupIP(ServerAddr))
	//return

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("initializing server")

	handler := NewHTTPImageProxyHandler(logger)

	server := &graceful.Server{
		Timeout: time.Second,
		Server: &http.Server{
			Addr:    ServerAddr,
			Handler: MiddlewareSecurity(handler, logger),
		},
		Logger: logger,
		ShutdownInitiated: func() {
			handler.Stop()
		},
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Println("starting server")

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger.Fatal(err)
		}
	}

	logger.Println("server stopped")
}

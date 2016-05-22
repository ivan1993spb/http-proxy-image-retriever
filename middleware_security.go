package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
)

const (
	URL_PATH         = "/"
	QUERY_PARAM_NAME = "url"
)

func MiddlewareSecurity(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("checking request")

		if r.Method != http.MethodGet {
			logger.Println("request method not allowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != URL_PATH {
			logger.Println("invalid path")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		URL := r.FormValue(QUERY_PARAM_NAME)

		if URL == "" {
			logger.Println("passed empty url")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		parsedUrl, err := url.Parse(URL)
		if err != nil {
			logger.Println("invalid url")
			return
		}

		logger.Printf("%#v\n", parsedUrl)

		if err = checkUrl(parsedUrl); err != nil {
			logger.Println("url problem:", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkUrl(URL *url.URL) error {
	host, port, err := net.SplitHostPort(URL.Host)
	if err != nil {
		return fmt.Errorf("checking url error: %s", err)
	}

	if host == "" {
		return errors.New("empty host")
	}

	if port == "" {
		if URL.Scheme == "http" {
			port = "80"
		} else if URL.Scheme == "https" {
			port = "443"
		} else {
			return errors.New("unknown port")
		}
	}

	if isHostPortOfThisServer(host, port) {
		return errors.New("ddos!")
	}

	return nil

}

func isHostPortOfThisServer(host, port string) bool {
	serverHost, serverPort, err := net.SplitHostPort(ServerAddr)
	if err != nil {
		// server cannot run with invalid addr
		return false
	}

	if serverPort == "" {
		serverPort = "80"
	}

	if serverPort != port {
		return false
	}

	serverAddrs, err := net.LookupAddr(serverHost)
	if err != nil {
		return false
	}

	addrs, err := net.LookupAddr(host)
	if err != nil {
		return false
	}

	log.Println(addrs)
	log.Println(serverAddrs)

	for _, addr := range addrs {
		for _, serverAddr := range serverAddrs {
			if addr == serverAddr {
				return true
			}
		}
	}

	return false
}

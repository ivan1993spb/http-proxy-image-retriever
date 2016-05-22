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
	fmt.Println(1)
	if err != nil {
		// server cannot run with invalid addr
		return false
	}
	fmt.Println(2)

	if serverPort == "" {
		serverPort = "80"
	}
	fmt.Println(3)

	if serverPort != port {
		return false
	}
	fmt.Println(4)

	fmt.Println(serverHost)

	var serverAddrs []string
	if serverHost == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return false
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					serverAddrs = append(serverAddrs, ipnet.IP.String())
				}
				if ipnet.IP.To16() != nil {
					serverAddrs = append(serverAddrs, ipnet.IP.String())
				}
			}

		}
	} else {
		serverAddrs, err = net.LookupHost(serverHost)
		if err != nil {
			return false
		}
	}
	fmt.Println(5)

	fmt.Println(host)
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Println(err)
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

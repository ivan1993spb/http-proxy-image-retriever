package main

import (
	"log"
	"net/http"
	"net/url"
)

func MiddlewareCheckRequest(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("accepted request")
		defer log.Println("finished request handling")
		logger.Println("checking request")

		if r.Method != http.MethodGet {
			logger.Println("request method not allowed")
			HTTPErrorHTML(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/" {
			logger.Println("unknown path")
			HTTPErrorHTML(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		URL := r.FormValue("url")

		if URL == "" {
			logger.Println("passed empty url")
			HTTPErrorHTML(w, "empty url param", http.StatusOK)
			return
		}

		parsedUrl, err := url.Parse(URL)
		if err != nil {
			logger.Println("invalid url:", err)
			HTTPErrorHTML(w, "invalid url", http.StatusOK)
			return
		}

		if parsedUrl.Host == "" {
			logger.Println("invalid url: empty host")
			HTTPErrorHTML(w, "invalid url", http.StatusOK)
			return
		}

		logger.Println("request is trusted")

		next.ServeHTTP(w, r)
	})
}

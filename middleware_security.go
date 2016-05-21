package main

import (
	"log"
	"net/http"
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

		next.ServeHTTP(w, r)
	})
}

package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go-bindata -o bindata_test.go test/index/ test/path/to/imgs/

func TestNewImageProxyHandler(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/index/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(MustAsset("test/index/index.html"))
	})
	mux.HandleFunc("/index/second.gif", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/gif")
		w.Write(MustAsset("test/index/second.gif"))
	})
	mux.HandleFunc("/index/sixth.jpg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(MustAsset("test/index/sixth.jpg"))
	})
	mux.HandleFunc("/index/third.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(MustAsset("test/index/third.svg"))
	})
	mux.HandleFunc("/path/to/imgs/fifth.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("test/path/to/imgs/fifth.png"))
	})
	mux.HandleFunc("/path/to/imgs/first.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("test/path/to/imgs/first.png"))
	})
	mux.HandleFunc("/path/to/imgs/fourth.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("test/path/to/imgs/fourth.png"))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	testURL := "http://localhost:8888?url=" + testServer.URL + "/index/index.html"
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		assert.FailNow(t, "cannot create test request", err)
	}
	w := httptest.NewRecorder()

	handler := NewImageProxyHandler(log.New(os.Stdout, "", log.LstdFlags))
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "invalud status code")

	// Save generated page to file if it possible
	if f, err := os.Create("test_result.html"); err == nil {
		io.Copy(f, w.Body)
		f.Close()
	}
}

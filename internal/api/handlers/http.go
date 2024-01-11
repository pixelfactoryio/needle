// Package handlers provides http handlers.
package handlers

import (
	"fmt"
	"net/http"
)

// Pixel 1x1 transparent pixel
var Pixel = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33,
	249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59,
}

type httpHandler struct{}

type caHandler struct {
	caFile string
}

// NewDefaultHandler create NewDefaultHandler
func NewDefaultHandler() http.Handler {
	return &httpHandler{}
}

// ServeHTTP respond with 1x1 transparent gif
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", fmt.Sprint(len(Pixel)))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(Pixel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// NewCAHandler create NewCAHandler
func NewCAHandler(ca string) http.Handler {
	return &caHandler{caFile: ca}
}

// ServeHTTP respond with public CA certificat
func (h *caHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-x509-ca-cert")
	http.ServeFile(w, r, h.caFile)
}

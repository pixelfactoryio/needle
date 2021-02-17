package handlers

import (
	"fmt"
	"net/http"
)

// 1x1 transparent pixel
var pixel = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33,
	249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59,
}

// HTTPHandlerService interface
type HTTPHandlerService interface {
	GetPixel(http.ResponseWriter, *http.Request)
	GetRootCA(http.ResponseWriter, *http.Request)
}

type HTTPHandler struct {
	caFile string
}

// NewPixelHandler create PixelHandler
func NewPixelHandler(caFile string) *HTTPHandler {
	return &HTTPHandler{caFile: caFile}
}

// GetPixel respond with 1x1 transparent gif
func (h *HTTPHandler) GetPixel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", fmt.Sprint(len(pixel)))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
	w.Write(pixel)
}

func (h *HTTPHandler) GetRootCA(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-x509-ca-cert")
	http.ServeFile(w, r, h.caFile)
}

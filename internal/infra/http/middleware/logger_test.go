package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestResponseWriterWriteHeaderOnce(t *testing.T) {
	is := require.New(t)

	recorder := httptest.NewRecorder()
	wrapped := wrapResponseWriter(recorder)

	wrapped.WriteHeader(http.StatusCreated)
	wrapped.WriteHeader(http.StatusBadRequest)

	is.Equal(http.StatusCreated, wrapped.Status())
	is.Equal(http.StatusCreated, recorder.Code)
}

func TestLoggingMiddlewareSuccess(t *testing.T) {
	is := require.New(t)
	logger := log.New()

	middleware := Logging(logger)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.RemoteAddr = "127.0.0.1:8080"

	handler.ServeHTTP(recorder, request)

	is.Equal(http.StatusNoContent, recorder.Code)
}

func TestLoggingMiddlewareInvalidRemoteAddr(t *testing.T) {
	is := require.New(t)
	logger := log.New()

	middleware := Logging(logger)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.RemoteAddr = "invalid"

	handler.ServeHTTP(recorder, request)

	is.Equal(http.StatusOK, recorder.Code)
}

func TestLoggingMiddlewareRecover(t *testing.T) {
	is := require.New(t)
	logger := log.New()

	middleware := Logging(logger)
	handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom")
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	handler.ServeHTTP(recorder, request)

	is.Equal(http.StatusInternalServerError, recorder.Code)
}

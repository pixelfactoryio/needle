// Package middleware provides HTTP middleware functions for logging and observability.
package middleware

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// Logging logs the incoming HTTP request & its duration.
func Logging(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("Internal Server Error", fields.Any("error", err))
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			ip, portStr, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				logger.Error("Unable to parse remote address", fields.String("remote_addr", r.RemoteAddr), fields.Error(err))
			}

			port, err := strconv.Atoi(portStr)
			if err != nil {
				logger.Error("Unable to parse remote port", fields.String("remote_port", portStr), fields.Error(err))
			}

			logger.Debug(
				"Request",
				fields.HTTPRequest(r),
				fields.Source(ip, port),
				fields.URL(r.URL),
				fields.UserAgent(r.UserAgent()),
				fields.Duration("duration", time.Since(start)),
				fields.Int("http.response.status_code", wrapped.status),
			)
		}

		return http.HandlerFunc(fn)
	}
}

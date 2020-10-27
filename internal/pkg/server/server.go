package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

// Handle OS Signals
var stopCh = setupSignalHandler()

// Config holds server config
type Config struct {
	Port                      string        `mapstructure:"http-port"`
	HTTPServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HTTPServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
}

// Server holds server
type Server struct {
	name      string
	router    http.Handler
	config    *Config
	logger    log.Logger
	tlsConfig *tls.Config
}

// Option type
type Option func(*Server)

// WithName set server name
func WithName(n string) Option {
	return func(s *Server) {
		s.name = n
	}
}

// WithRouter set server http Handler
func WithRouter(r http.Handler) Option {
	return func(s *Server) {
		s.router = r
	}
}

// WithConfig set server config
func WithConfig(c *Config) Option {
	return func(s *Server) {
		s.config = c
	}
}

// WithLogger set server logger
func WithLogger(l log.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// WithTLSConfig set server tls configuration
func WithTLSConfig(c *tls.Config) Option {
	return func(s *Server) {
		s.tlsConfig = c
	}
}

// NewServer create new Server with default values
func NewServer(opts ...Option) *Server {
	// setup default server
	srv := &Server{
		name:   "default",
		router: http.NewServeMux(),
		config: &Config{
			Port: "443",
		},
	}

	for _, opt := range opts {
		opt(srv)
	}

	// setup default logger
	if srv.logger == nil {
		srv.logger = log.New()
		srv.logger.Info("Using default logger")
	}

	// setup default server timeout
	if srv.config.HTTPServerTimeout == 0 {
		srv.config.HTTPServerTimeout = 60 * time.Second
		srv.logger.Debug("Using default HTTPServerTimeout", fields.String("server-timeout", srv.config.HTTPServerTimeout.String()))
	}

	// setup default server shutdown timeout
	if srv.config.HTTPServerShutdownTimeout == 0 {
		srv.config.HTTPServerShutdownTimeout = 5 * time.Second
		srv.logger.Debug("Using default HTTPServerShutdownTimeout", fields.String("server-timeout", srv.config.HTTPServerShutdownTimeout.String()))
	}

	return srv
}

// ListenAndServe start server
func (s *Server) ListenAndServe() {
	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		WriteTimeout: s.config.HTTPServerTimeout,
		ReadTimeout:  s.config.HTTPServerTimeout,
		IdleTimeout:  2 * s.config.HTTPServerTimeout,
	}

	// Create listener
	var ln net.Listener
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))
	if err != nil {
		s.logger.Error("Unable to create net.Listener", fields.Error(err))
	}

	if s.tlsConfig != nil {
		s.tlsConfig.NextProtos = append(s.tlsConfig.NextProtos, "h2")
		s.tlsConfig.NextProtos = append(s.tlsConfig.NextProtos, "http/1.1")
		ln = tls.NewListener(ln, s.tlsConfig)
	}
	defer ln.Close()

	// run server in background
	go func() {
		s.logger.Info("Starting server")
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			s.logger.Error("Server crashed", fields.Error(err))
		}
	}()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HTTPServerShutdownTimeout)
	defer cancel()

	s.logger.Info("Shutting down server")

	// attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("Server graceful shutdown failed", fields.Error(err))
	} else {
		s.logger.Info("Server stopped")
	}
}

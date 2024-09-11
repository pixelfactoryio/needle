package cmd

import (
	"crypto/tls"
	"time"

	"github.com/spf13/cobra"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
	"go.pixelfactory.io/pkg/server"
	"go.pixelfactory.io/pkg/version"

	"go.pixelfactory.io/needle/internal/app/factory"
	"go.pixelfactory.io/needle/internal/app/pki"
	"go.pixelfactory.io/needle/internal/infra/boltdb"
	"go.pixelfactory.io/needle/internal/infra/coredns"
	"go.pixelfactory.io/needle/internal/infra/http"
	"go.pixelfactory.io/needle/internal/infra/http/handlers"
)

var (
	logLevel                  string
	caFile                    string
	caKeyFile                 string
	dbFile                    string
	httpPort                  string
	httpsPort                 string
	httpServerTimeout         time.Duration
	httpServerShutdownTimeout time.Duration
	corednsEnabled            bool
	corednsPort               int
	corednsHostsFile          string
	corednsUpstreams          []string
	corednsCoreFile           string
)

var needleCmd = &cobra.Command{
	Use:   "needle",
	Short: "needle",
	RunE:  start,
}

// NewNeedleCmd create new needleCmd.
func NewNeedleCmd() (*cobra.Command, error) {
	needleCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")
	if err := bindFlag("log-level"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&caFile, "ca", "data/certs/root-ca.crt", "Root CA Certificate path")
	if err := bindFlag("ca"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&caKeyFile, "ca-key", "data/certs/root-ca.key", "Root CA Key path")
	if err := bindFlag("ca-key"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&dbFile, "db-file", "data/cache.db", "Cache DB path")
	if err := bindFlag("db-file"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&httpPort, "http-port", "80", "HTTP port")
	if err := bindFlag("http-port"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&httpsPort, "https-port", "443", "HTTPS port")
	if err := bindFlag("https-port"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().DurationVar(&httpServerTimeout, "server-timeout", 60*time.Second, "Server timeout")
	if err := bindFlag("server-timeout"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().DurationVar(
		&httpServerShutdownTimeout, "server-shutdown-timeout", 5*time.Second, "Server shutdown timeout")
	if err := bindFlag("server-shutdown-timeout"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().BoolVar(&corednsEnabled, "coredns", false, "Enable embedded CoreDNS")
	if err := bindFlag("coredns"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().IntVar(&corednsPort, "coredns-port", 53, "CoreDNS port")
	if err := bindFlag("coredns-port"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&corednsHostsFile, "coredns-hosts-file", "data/hosts", "Hosts file path")
	if err := bindFlag("coredns-hosts-file"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringSliceVar(
		&corednsUpstreams, "coredns-upstreams", []string{"1.1.1.1", "8.8.8.8"}, "Upstream DNS servers")
	if err := bindFlag("coredns-upstreams"); err != nil {
		return nil, err
	}

	needleCmd.PersistentFlags().StringVar(&corednsCoreFile, "coredns-corefile", "data/Corefile", "Corefile file path")
	if err := bindFlag("coredns-corefile"); err != nil {
		return nil, err
	}

	return needleCmd, nil
}

func start(_ *cobra.Command, _ []string) error {
	// Setup logger
	logger := log.New(log.WithLevel(logLevel))
	logger = logger.With(fields.Service("needle", version.REVISION))
	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Error("an error occurred while running logger.Sync()", fields.Error(err))
		}
	}()

	logger.Debug("Needle Start")
	logger.Debug(
		"Needle Configuration",
		fields.String("caFile", caFile),
		fields.String("keyFile", caKeyFile),
		fields.String("dbFile", dbFile),
		fields.String("http-port", httpPort),
		fields.String("https-port", httpsPort),
		fields.String("server-timeout", httpServerTimeout.String()),
		fields.String("server-shutdown-timeout", httpServerShutdownTimeout.String()),
	)

	if corednsEnabled {
		dnsServer := coredns.NewCoreDNSServer(
			coredns.WithLogger(logger),
			coredns.WithPort(corednsPort),
			coredns.WithHostsFile(corednsHostsFile),
			coredns.WithUpstreams(corednsUpstreams),
			coredns.WithCoreFile(corednsCoreFile),
		)

		// Start CoreDNS Server
		go func() {
			if err := dnsServer.Run(); err != nil {
				logger.Error("failed to run CoreDNS server", fields.Error(err))
			}
		}()

		logger.Debug("CoreDNS Server started")
		logger.Debug(
			"CoreDNS Configuration",
			fields.Int("port", corednsPort),
			fields.String("hostsfile", corednsHostsFile),
			fields.Strings("upstreams", corednsUpstreams),
			fields.String("corefile", corednsCoreFile),
		)
	}

	// Setup tls certificate service
	rootCA, err := tls.LoadX509KeyPair(caFile, caKeyFile)
	if err != nil {
		return err
	}

	// Setup BoltDB repository
	client, err := newStormClient(dbFile)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close()
		if err != nil {
			logger.Error("an error occurred while closing *storm.DB client", fields.Error(err))
		}
	}()

	// Setup PKI service
	pkiSvc := pki.New(
		boltdb.New(client),
		factory.New(rootCA),
	)

	// Setup certificate handler and tls.Config
	certHandler := handlers.NewTLSHandler(logger, pkiSvc)
	tlsConfig := &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: certHandler,
		NextProtos:     []string{"h2", "http/1.1"},
	}

	// Setup http handler
	routes := []http.Route{
		{
			Path:    "/install-root-ca",
			Handler: handlers.NewCAHandler(caFile),
		},
		{
			Path:    "/",
			Handler: handlers.NewDefaultHandler(),
		},
	}

	router := http.NewRouter(logger, routes...)

	tlsSrv, err := server.New(
		server.WithName("needle-tls"),
		server.WithLogger(logger),
		server.WithRouter(router),
		server.WithPort(httpsPort),
		server.WithHTTPServerTimeout(httpServerTimeout),
		server.WithHTTPServerShutdownTimeout(httpServerShutdownTimeout),
		server.WithTLSConfig(tlsConfig),
	)
	if err != nil {
		return err
	}

	// Start TLS Server
	go func() {
		err := tlsSrv.ListenAndServe()
		if err != nil {
			logger.Error("TLS server failed to start", fields.Error(err))
		}
	}()

	httpSrv, err := server.New(
		server.WithName("needle-http"),
		server.WithLogger(logger),
		server.WithRouter(router),
		server.WithPort(httpPort),
		server.WithHTTPServerTimeout(httpServerTimeout),
		server.WithHTTPServerShutdownTimeout(httpServerShutdownTimeout),
	)
	if err != nil {
		return err
	}

	// Start HTTP Server
	if err := httpSrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

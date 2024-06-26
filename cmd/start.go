package cmd

import (
	"crypto/tls"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/spf13/cobra"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
	"go.pixelfactory.io/pkg/server"
	"go.pixelfactory.io/pkg/version"

	"go.pixelfactory.io/needle/internal/api"
	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/needle/internal/pkg/coredns"
	"go.pixelfactory.io/needle/internal/repository/boltdb"
	"go.pixelfactory.io/needle/internal/services/pki"
)

var (
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
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start needle",
	RunE:  start,
}

// NewStartCmd command topic
// nolint
func NewStartCmd() (*cobra.Command, error) {
	startCmd.PersistentFlags().StringVar(&caFile, "ca", "data/certs/root-ca.crt", "Root CA Certificate path")
	if err := bindFlag("ca"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringVar(&caKeyFile, "ca-key", "data/certs/root-ca.key", "Root CA Key path")
	if err := bindFlag("ca-key"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringVar(&dbFile, "db-file", "data/cache.db", "Cache DB path")
	if err := bindFlag("db-file"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringVar(&httpPort, "http-port", "80", "HTTP port")
	if err := bindFlag("http-port"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringVar(&httpsPort, "https-port", "443", "HTTPS port")
	if err := bindFlag("https-port"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().DurationVar(&httpServerTimeout, "server-timeout", 60*time.Second, "Server timeout")
	if err := bindFlag("server-timeout"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().DurationVar(
		&httpServerShutdownTimeout, "server-shutdown-timeout", 5*time.Second, "Server shutdown timeout")
	if err := bindFlag("server-shutdown-timeout"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().BoolVar(&corednsEnabled, "coredns", false, "Enable embedded CoreDNS")
	if err := bindFlag("coredns"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().IntVar(&corednsPort, "coredns-port", 53, "CoreDNS port")
	if err := bindFlag("coredns-port"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringVar(&corednsHostsFile, "coredns-hosts-file", "data/hosts", "Hosts file path")
	if err := bindFlag("coredns-hosts-file"); err != nil {
		return nil, err
	}

	startCmd.PersistentFlags().StringSliceVar(
		&corednsUpstreams, "coredns-upstreams", []string{"1.1.1.1", "8.8.8.8"}, "Upstream DNS servers")
	if err := bindFlag("coredns-upstreams"); err != nil {
		return nil, err
	}

	return startCmd, nil
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

	logger.Debug("Application Start")
	logger.Debug(
		"Application Configuration",
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
			coredns.WithPort(corednsPort),
			coredns.WithHostsFile(corednsHostsFile),
			coredns.WithUpstreams(corednsUpstreams),
		)
		go dnsServer.Run()
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

	// Setup service
	service := pki.NewCertificateService(
		boltdb.NewBoltRepository(client),
		pki.NewFactory(rootCA),
	)

	// Setup certificate handler and tls.Config
	certHandler := handlers.NewTLSHandler(logger, service)
	tlsConfig := &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: certHandler,
	}

	// Setup http handler
	routes := []api.Route{
		{
			Path:    "/install-root-ca",
			Handler: handlers.NewCAHandler(caFile),
		},
		{
			Path:    "/",
			Handler: handlers.NewDefaultHandler(),
		},
	}

	router := api.NewRouter(logger, routes...)

	tlsSrv, err := server.NewServer(
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
	go tlsSrv.ListenAndServe()

	httpSrv, _ := server.NewServer(
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
	httpSrv.ListenAndServe()

	return nil
}

func newStormClient(dbFile string) (*storm.DB, error) {
	client, err := storm.Open(dbFile)
	if err != nil {
		return nil, err
	}
	return client, nil
}

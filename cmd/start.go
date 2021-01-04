package cmd

import (
	"crypto/tls"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/spf13/cobra"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"

	"go.pixelfactory.io/needle/internal/api"
	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/needle/internal/services/pki"

	"go.pixelfactory.io/needle/internal/pkg/coredns"
	"go.pixelfactory.io/needle/internal/pkg/server"
	"go.pixelfactory.io/needle/internal/pkg/version"
	"go.pixelfactory.io/needle/internal/repository/boltdb"
)

var rootCA tls.Certificate
var caFile string
var caKeyFile string
var dbFile string
var cacheDB *storm.DB
var httpPort string
var httpsPort string
var httpServerTimeout time.Duration
var httpServerShutdownTimeout time.Duration
var corednsEnabled bool
var corednsPort int
var corednsHostsFile string
var corednsUpstreams []string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start needle",
	RunE:  start,
}

// NewStartCmd command topic
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

	startCmd.PersistentFlags().DurationVar(&httpServerShutdownTimeout, "server-shutdown-timeout", 5*time.Second, "Server shutdown timeout")
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

	startCmd.PersistentFlags().StringSliceVar(&corednsUpstreams, "coredns-upstreams", []string{"1.1.1.1", "8.8.8.8"}, "Upstream DNS servers")
	if err := bindFlag("coredns-upstreams"); err != nil {
		return nil, err
	}

	return startCmd, nil
}

func start(c *cobra.Command, args []string) error {
	// Setup logger
	logger := log.New(log.WithLevel(logLevel))
	logger = logger.With(fields.Service("needle", version.REVISION))
	defer logger.Sync()

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
	defer client.Close()

	// Setup service
	service := pki.NewCertificateService(
		boltdb.NewBoltRepository(client),
		pki.NewFactory(rootCA),
	)

	// Setup certificate handler and tls.Config
	certHandler := handlers.NewCertificateHandler(logger, service)
	tlsConfig := &tls.Config{GetCertificate: certHandler.Get}

	// Setup http handler
	router := api.NewRouter(logger, caFile)

	tlsSrv := server.NewServer(
		server.WithLogger(logger),
		server.WithRouter(router),
		server.WithConfig(&server.Config{
			Port:                      httpsPort,
			HTTPServerTimeout:         httpServerTimeout,
			HTTPServerShutdownTimeout: httpServerShutdownTimeout,
		}),
		server.WithTLSConfig(tlsConfig),
	)
	go tlsSrv.ListenAndServe()

	httpSrv := server.NewServer(
		server.WithLogger(logger),
		server.WithRouter(router),
		server.WithConfig(&server.Config{
			Port:                      httpPort,
			HTTPServerTimeout:         httpServerTimeout,
			HTTPServerShutdownTimeout: httpServerShutdownTimeout,
		}),
	)
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

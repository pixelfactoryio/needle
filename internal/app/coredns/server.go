// Package coredns provides a DNS server implementation using CoreDNS.
package coredns

import (
	"bytes"
	"text/template"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

// DNSServer holds dns server.
type DNSServer struct {
	name      string
	port      int
	hostsfile string
	upsteams  []string
	logger    log.Logger
}

// Option type.
type Option func(*DNSServer)

// WithLogger set server logger.
func WithLogger(l log.Logger) Option {
	return func(s *DNSServer) {
		s.logger = l
	}
}

// WithPort set server port.
func WithPort(p int) Option {
	return func(s *DNSServer) {
		s.port = p
	}
}

// WithHostsFile set hosts file.
func WithHostsFile(h string) Option {
	return func(s *DNSServer) {
		s.hostsfile = h
	}
}

// WithUpstreams set upstream dns servers.
func WithUpstreams(u []string) Option {
	return func(s *DNSServer) {
		s.upsteams = u
	}
}

// NewCoreDNSServer create new DNSServer with default values.
func NewCoreDNSServer(opts ...Option) *DNSServer {
	srv := &DNSServer{
		name:      "coredns",
		port:      53,
		hostsfile: "hosts",
		upsteams:  []string{"/etc/resolv.conf"},
	}

	for _, opt := range opts {
		opt(srv)
	}

	// setup default logger
	if srv.logger == nil {
		srv.logger = log.New()
		srv.logger.Info("Using default logger")
	}

	caddy.Quiet = true // don't show init stuff from caddy
	dnsserver.Quiet = true
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(srv.defaultLoader))

	return srv
}

// Run CoreDNS.
func (s *DNSServer) Run() {
	corefile, err := caddy.LoadCaddyfile("dns")
	if err != nil {
		s.logger.Error("DNS Server crashed", fields.Error(err))
	}

	instance, err := caddy.Start(corefile)
	if err != nil {
		s.logger.Error("DNS Server crashed", fields.Error(err))
	}
	instance.Wait()
}

// defaultLoader loads the CoreDNS configuration.
func (s *DNSServer) defaultLoader(serverType string) (caddy.Input, error) {
	return caddy.CaddyfileInput{
		Contents:       s.renderCorefile(),
		ServerTypeName: serverType,
	}, nil
}

func (s *DNSServer) renderCorefile() []byte {
	corefileTpl := `
	.:{{.Port}} {
		hosts {{.Hosts}} {
			fallthrough
		}
		forward . {{range $server := .Upstreams}} {{$server}} {{end}}
		cache
		loop
	}`

	values := struct {
		Port      int
		Hosts     string
		Upstreams []string
	}{
		Port:      s.port,
		Hosts:     s.hostsfile,
		Upstreams: s.upsteams,
	}

	tmpl, err := template.New("corefile").Parse(corefileTpl)
	if err != nil {
		s.logger.Error("DNS Server crashed", fields.Error(err))
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, values)
	if err != nil {
		s.logger.Error("DNS Server crashed", fields.Error(err))
	}

	return tpl.Bytes()
}

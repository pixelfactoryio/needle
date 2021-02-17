package handlers

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"

	"go.pixelfactory.io/needle/internal/services/pki"
)

// CertificateHandlerFunc returns a Certificate based on the given ClientHelloInfo
type CertificateHandlerFunc func(*tls.ClientHelloInfo) (*tls.Certificate, error)

// NewTLSHandler creates tlsHandler
func NewTLSHandler(logger log.Logger, certificateService pki.CertificateService) CertificateHandlerFunc {
	return func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
		logger.Debug("Getting certificate",
			fields.String("ServerName", helloInfo.ServerName), fields.String("LocalAddr", helloInfo.Conn.LocalAddr().String()),
		)

		name := getHostIP(helloInfo.Conn.LocalAddr().String())
		if len(helloInfo.ServerName) > 0 {
			name = helloInfo.ServerName
		}

		certificate, err := certificateService.GetOrCreate(name)
		if err != nil {
			err := errors.Wrap(err, "api.CertificateHandler.Get")
			logger.Error("Unable to find certificate", fields.String("CommonName", name), fields.Error(err))
			return nil, err
		}

		tlsCert, err := tls.X509KeyPair(certificate.CertPEM, certificate.KeyPEM)
		if err != nil {
			err := errors.Wrap(err, "api.CertificateHandler.Get")
			logger.Error("Error creating certificate", fields.String("CommonName", name), fields.Error(err))
			return nil, err
		}

		return &tlsCert, nil
	}
}

func getHostIP(localAddr string) string {
	u, err := url.Parse(fmt.Sprintf("https://%s", localAddr))
	if err != nil {
		return ""
	}
	return u.Hostname()
}

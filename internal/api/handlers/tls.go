package handlers

import (
	"crypto/tls"

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
		logger.Debug("Getting certificate", fields.String("ServerName", helloInfo.ServerName))

		name := "default-needle-certificate"
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

package handlers

import (
	"crypto/tls"

	"github.com/pkg/errors"
	"go.pixelfactory.io/needle/internal/app/pki"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

type PKIService interface {
	GetOrCreate(name string) (*pki.InternalCert, error)
}

// CertificateHandlerFunc returns a Certificate based on the given ClientHelloInfo.
type CertificateHandlerFunc func(*tls.ClientHelloInfo) (*tls.Certificate, error)

// NewTLSHandler creates tlsHandler.
func NewTLSHandler(logger log.Logger, pkiSvc PKIService) CertificateHandlerFunc {
	return func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
		logger.Debug("Getting certificate", fields.String("ServerName", helloInfo.ServerName))

		name := "default-needle-certificate"
		if helloInfo.ServerName != "" {
			name = helloInfo.ServerName
		}

		certificate, err := pkiSvc.GetOrCreate(name)
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

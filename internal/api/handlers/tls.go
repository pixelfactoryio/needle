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

// TLSHandlerService interface
type TLSHandlerService interface {
	Get(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error)
}

type TLSHandler struct {
	certificateService pki.CertificateService
	logger             log.Logger
}

// NewTLSHandler creates tlsHandler
func NewTLSHandler(logger log.Logger, certificateService pki.CertificateService) *TLSHandler {
	return &TLSHandler{
		certificateService: certificateService,
		logger:             logger,
	}
}

// Get tls.Certificate from given name
func (h *TLSHandler) Get(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	h.logger.Debug("Getting certificate",
		fields.String("ServerName", helloInfo.ServerName), fields.String("LocalAddr", helloInfo.Conn.LocalAddr().String()),
	)

	name := getHostIP(helloInfo.Conn.LocalAddr().String())
	if len(helloInfo.ServerName) > 0 {
		name = helloInfo.ServerName
	}

	certificate, err := h.certificateService.GetOrCreate(name)
	if err != nil {
		err := errors.Wrap(err, "api.CertificateHandler.Get")
		h.logger.Error("Unable to find certificate", fields.String("CommonName", name), fields.Error(err))
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(certificate.CertPEM, certificate.KeyPEM)
	if err != nil {
		err := errors.Wrap(err, "api.CertificateHandler.Get")
		h.logger.Error("Error creating certificate", fields.String("CommonName", name), fields.Error(err))
		return nil, err
	}

	return &tlsCert, nil
}

func getHostIP(localAddr string) string {
	u, err := url.Parse(fmt.Sprintf("https://%s", localAddr))
	if err != nil {
		return ""
	}
	return u.Hostname()
}

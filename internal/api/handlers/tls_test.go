package handlers_test

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/needle/internal/services/pki"
	"go.pixelfactory.io/needle/mocks/pkimock"
	"go.pixelfactory.io/needle/testdata"
	"go.pixelfactory.io/pkg/observability/log"
)

func Test_TLSHandler(t *testing.T) {
	t.Parallel()
	is := require.New(t)
	logger := log.New()

	rootCA, testCert := testdata.Setup(t)
	x509CACert, err := x509.ParseCertificate(rootCA.Certificate[0])
	if err != nil {
		t.Error("Error:", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(x509CACert)

	svc := &pkimock.CertificateService{}
	tlsHandler := handlers.NewTLSHandler(logger, svc)

	is.NotEmpty(tlsHandler)
	is.NotEmpty(tlsHandler)
	is.IsType((handlers.CertificateHandlerFunc)(nil), tlsHandler)

	t.Run("Create certificate", func(t *testing.T) {
		svc.On("GetOrCreate", "test.needle.local").Return(testCert, nil).Once()

		tlsCert, err := tlsHandler(&tls.ClientHelloInfo{ServerName: "test.needle.local"})
		if err != nil {
			t.Error(err)
		}

		is.NotEmpty(tlsCert)

		x509tlsCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
		if err != nil {
			t.Error("Error:", err)
		}

		t.Run("Verify certificate", func(t *testing.T) {
			_, err = x509tlsCert.Verify(x509.VerifyOptions{DNSName: "test.needle.local", Roots: roots})
			is.NoError(err)
		})
	})

	t.Run("Create certificate error", func(t *testing.T) {
		svc.On("GetOrCreate", "test.needle.local").Return(nil, errors.New("unable to create certificate")).Once()

		tlsCert, err := tlsHandler(&tls.ClientHelloInfo{ServerName: "test.needle.local"})
		is.Error(err)
		is.Empty(tlsCert)
	})

	t.Run("Create certificate error empty certificate", func(t *testing.T) {
		svc.On("GetOrCreate", "test.needle.local").Return(&pki.Certificate{}, nil).Once()

		tlsCert, err := tlsHandler(&tls.ClientHelloInfo{ServerName: "test.needle.local"})
		is.Error(err)
		is.Empty(tlsCert)
	})
}

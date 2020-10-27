package pki_test

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/services/pki"
)

func Test_NewFactory(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	rootCA, _ := setup(t)

	certFactory := pki.NewFactory(rootCA)
	is.NotEmpty(certFactory)
	is.Implements((*pki.CertificateFactory)(nil), certFactory)
}

func Test_Create(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	rootCA, _ := setup(t)
	x509CACert, err := x509.ParseCertificate(rootCA.Certificate[0])
	if err != nil {
		t.Error("Error:", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(x509CACert)

	certFactory := pki.NewFactory(rootCA)

	t.Run("Create certificate", func(t *testing.T) {
		// create certificate
		cert, err := certFactory.Create("test.needle.local")
		is.NoError(err)
		is.Equal(cert.Name, "test.needle.local")

		// convert to tls.Certificate
		tlsCert, err := tls.X509KeyPair(cert.CertPEM, cert.KeyPEM)
		if err != nil {
			t.Error("Error:", err)
		}

		// convert to x509.Certificate
		x509tlsCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
		if err != nil {
			t.Error("Error:", err)
		}

		t.Run("Verify certificate", func(t *testing.T) {
			_, err = x509tlsCert.Verify(x509.VerifyOptions{DNSName: "test.needle.local", Roots: roots})
			is.NoError(err)
		})
	})
}

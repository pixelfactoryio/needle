package pki_test

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/services/pki"
	"go.pixelfactory.io/needle/testdata"
)

func Test_NewFactory(t *testing.T) {
	is := require.New(t)

	rootCA, _ := testdata.Setup(t)

	certFactory := pki.NewFactory(rootCA)
	is.NotEmpty(certFactory)
	is.Implements((*pki.CertificateFactory)(nil), certFactory)
}

func Test_Create(t *testing.T) {
	is := require.New(t)

	rootCA, _ := testdata.Setup(t)
	x509CACert, err := x509.ParseCertificate(rootCA.Certificate[0])
	is.NoError(err)

	roots := x509.NewCertPool()
	roots.AddCert(x509CACert)

	certFactory := pki.NewFactory(rootCA)

	t.Run("Create certificate", func(_ *testing.T) {
		// create certificate
		cert, err := certFactory.Create("test.needle.local")
		is.NoError(err)
		is.Equal(cert.Name, "test.needle.local")

		// convert to tls.Certificate
		tlsCert, err := tls.X509KeyPair(cert.CertPEM, cert.KeyPEM)
		is.NoError(err)

		// convert to x509.Certificate
		x509tlsCert, err := x509.ParseCertificate(tlsCert.Certificate[0])
		is.NoError(err)

		_, err = x509tlsCert.Verify(x509.VerifyOptions{DNSName: "test.needle.local", Roots: roots})
		is.NoError(err)
	})

	t.Run("Create certificate IP", func(_ *testing.T) {
		// create certificate
		cert, err := certFactory.Create("192.168.1.1")
		is.NoError(err)
		is.Equal(cert.Name, "192.168.1.1")
	})
}

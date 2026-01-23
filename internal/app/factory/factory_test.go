package factory

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pixelfactory.io/needle/internal/app/pki"
	"go.pixelfactory.io/needle/testdata"
)

func Test_NewFactory(t *testing.T) {
	is := require.New(t)

	rootCA, _ := testdata.Setup(t)

	certFactory := New(rootCA)
	is.NotEmpty(certFactory)
	is.Implements((*pki.Factory)(nil), certFactory)
}

func Test_Create(t *testing.T) {
	is := require.New(t)

	rootCA, _ := testdata.Setup(t)
	x509CACert, parseErr := x509.ParseCertificate(rootCA.Certificate[0])
	is.NoError(parseErr)

	roots := x509.NewCertPool()
	roots.AddCert(x509CACert)

	certFactory := New(rootCA)

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

	t.Run("Create certificate invalid CA", func(_ *testing.T) {
		invalidCA := tls.Certificate{Certificate: [][]byte{[]byte("invalid-ca")}}
		invalidFactory := New(invalidCA)

		cert, err := invalidFactory.Create("test.needle.local")
		is.Error(err)
		is.Nil(cert)
	})

	t.Run("Create certificate signing error", func(_ *testing.T) {
		badRootCA, _ := testdata.Setup(t)
		badRootCA.PrivateKey = nil

		badFactory := New(badRootCA)
		cert, err := badFactory.Create("test.needle.local")
		is.Error(err)
		is.Nil(cert)
	})

	t.Run("Create certificate pem encode error", func(_ *testing.T) {
		originalEncode := pemEncode
		defer func() {
			pemEncode = originalEncode
		}()

		pemEncode = func(_ io.Writer, _ *pem.Block) error {
			return errors.New("pem error")
		}

		cert, err := certFactory.Create("test.needle.local")
		is.Error(err)
		is.Nil(cert)
	})

	t.Run("Create certificate pem encode key error", func(_ *testing.T) {
		originalEncode := pemEncode
		defer func() {
			pemEncode = originalEncode
		}()

		callCount := 0
		pemEncode = func(_ io.Writer, _ *pem.Block) error {
			callCount++
			if callCount == 2 {
				return errors.New("pem key error")
			}
			return nil
		}

		cert, err := certFactory.Create("test.needle.local")
		is.Error(err)
		is.Nil(cert)
	})
}

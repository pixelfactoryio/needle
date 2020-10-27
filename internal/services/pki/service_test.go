package pki_test

import (
	"crypto/tls"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/services/pki"
	"go.pixelfactory.io/needle/mocks/pkimock"
)

func setup(t *testing.T) (tls.Certificate, *pki.Certificate) {
	rootCA, err := tls.LoadX509KeyPair("testdata/certs/root-ca.crt", "testdata/certs/root-ca.key")
	if err != nil {
		t.Error("Unable to get rootCA", err)
	}

	certPEM, err := ioutil.ReadFile("testdata/certs/test.needle.local.crt")
	keyPEM, err := ioutil.ReadFile("testdata/certs/test.needle.local.key")

	testCert := pki.Certificate{
		Name:    "test.needle.local",
		CertPEM: certPEM,
		KeyPEM:  keyPEM,
	}
	return rootCA, &testCert
}

func Test_NewCertificateService(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	factory := &pkimock.CertificateFactory{}
	repo := &pkimock.CertificateRepository{}
	svc := pki.NewCertificateService(repo, factory)
	is.NotEmpty(svc)
	is.Implements((*pki.CertificateService)(nil), svc)
}

func Test_GetOrCreate(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	_, testCert := setup(t)
	repo := &pkimock.CertificateRepository{}
	factory := &pkimock.CertificateFactory{}
	svc := pki.NewCertificateService(repo, factory)

	t.Run("Create certificate", func(t *testing.T) {
		repo.On("Get", "test.needle.local").Return(nil, pki.ErrCertificateNotFound).Once()
		repo.On("Store", testCert).Return(nil).Once()
		factory.On("Create", "test.needle.local").Return(testCert, nil).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.NoError(err)
		is.NotEmpty(cert)
		repo.AssertExpectations(t)
	})

	t.Run("Get certificate", func(t *testing.T) {
		repo.On("Get", "test.needle.local").Return(testCert, nil).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.NoError(err)
		is.NotEmpty(cert.Name)
		is.NotEmpty(cert.CertPEM)
		is.NotEmpty(cert.KeyPEM)
		is.Equal(cert, testCert)
		repo.AssertExpectations(t)
	})
}

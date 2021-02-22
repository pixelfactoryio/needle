package pki_test

import (
	"crypto/tls"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
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
	if err != nil {
		t.Error("Unable to get certPEM", err)
	}

	keyPEM, err := ioutil.ReadFile("testdata/certs/test.needle.local.key")
	if err != nil {
		t.Error("Unable to get keyPEM", err)
	}

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
		factory.On("Create", "test.needle.local").Return(testCert, nil).Once()
		repo.On("Store", testCert).Return(nil).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.NoError(err)
		is.NotEmpty(cert)
		repo.AssertExpectations(t)
		factory.AssertExpectations(t)
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

	t.Run("Get certificate create error", func(t *testing.T) {
		repo.On("Get", "test.needle.local").Return(nil, pki.ErrCertificateNotFound).Once()
		factory.On("Create", "test.needle.local").Return(nil, errors.New("unable to create certificate")).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.Error(err)
		is.Empty(cert)
		repo.AssertExpectations(t)
		factory.AssertExpectations(t)
	})

	t.Run("Get certificate store error", func(t *testing.T) {
		repo.On("Get", "test.needle.local").Return(nil, pki.ErrCertificateNotFound).Once()
		factory.On("Create", "test.needle.local").Return(testCert, nil).Once()
		repo.On("Store", testCert).Return(errors.New("unable to store certificate")).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.Error(err)
		is.Empty(cert)
		repo.AssertExpectations(t)
		factory.AssertExpectations(t)
	})
}

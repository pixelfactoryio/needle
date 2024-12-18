package pki_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/app/pki"
	"go.pixelfactory.io/needle/internal/infra/http/handlers"
	mocks "go.pixelfactory.io/needle/mocks/pki"
	"go.pixelfactory.io/needle/testdata"
)

func Test_NewCertificateService(t *testing.T) {
	is := require.New(t)

	factory := &mocks.Factory{}
	repo := &mocks.Repository{}
	svc := pki.New(repo, factory)
	is.NotEmpty(svc)
	is.Implements((*handlers.PKIService)(nil), svc)
}

func Test_GetOrCreate(t *testing.T) {
	is := require.New(t)

	_, testCert := testdata.Setup(t)
	factory := &mocks.Factory{}
	repo := &mocks.Repository{}

	svc := pki.New(repo, factory)

	t.Run("Create certificate", func(_ *testing.T) {
		repo.On("Get", "test.needle.local").Return(nil, pki.ErrCertificateNotFound).Once()
		factory.On("Create", "test.needle.local").Return(testCert, nil).Once()
		repo.On("Store", testCert).Return(nil).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.NoError(err)
		is.NotEmpty(cert)
		repo.AssertExpectations(t)
		factory.AssertExpectations(t)
	})

	t.Run("Get certificate", func(_ *testing.T) {
		repo.On("Get", "test.needle.local").Return(testCert, nil).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.NoError(err)
		is.NotEmpty(cert.Name)
		is.NotEmpty(cert.CertPEM)
		is.NotEmpty(cert.KeyPEM)
		is.Equal(cert, testCert)
		repo.AssertExpectations(t)
	})

	t.Run("Get certificate create error", func(_ *testing.T) {
		repo.On("Get", "test.needle.local").Return(nil, pki.ErrCertificateNotFound).Once()
		factory.On("Create", "test.needle.local").Return(nil, errors.New("unable to create certificate")).Once()

		cert, err := svc.GetOrCreate("test.needle.local")
		is.Error(err)
		is.Empty(cert)
		repo.AssertExpectations(t)
		factory.AssertExpectations(t)
	})

	t.Run("Get certificate store error", func(_ *testing.T) {
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

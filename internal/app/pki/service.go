package pki

import (
	"github.com/pkg/errors"
)

// ErrCertificateNotFound unable to find certificate.
var ErrCertificateNotFound = errors.New("Certificate Not Found")

// Factory interface.
type Factory interface {
	Create(name string) (*InternalCert, error)
}

// Repository interface.
type Repository interface {
	Get(name string) (*InternalCert, error)
	Store(certificate *InternalCert) error
}

// Service represents a Certificate service.
type Service struct {
	certRepo    Repository
	certFactory Factory
}

// New create new service.
func New(certRepo Repository, certFactory Factory) *Service {
	return &Service{
		certRepo:    certRepo,
		certFactory: certFactory,
	}
}

// GetOrCreate retrives or create a certificat for the given name.
func (s *Service) GetOrCreate(name string) (*InternalCert, error) {
	var cert *InternalCert
	cert, err := s.certRepo.Get(name)
	if errors.Is(err, ErrCertificateNotFound) {
		// Create new Certificate
		cert, err = s.certFactory.Create(name)
		if err != nil {
			return nil, errors.Wrap(err, "pki.Service.GetOrCreate")
		}
		// Store Certificate
		err := s.certRepo.Store(cert)
		if err != nil {
			return nil, errors.Wrap(err, "pki.Service.GetOrCreate")
		}

	}
	return cert, nil
}

package pki

import (
	"github.com/pkg/errors"
)

// ErrCertificateNotFound unable to find certificate.
var ErrCertificateNotFound = errors.New("Certificate Not Found")

// CertificateService the necessary functionality required by service to perform operation.
type CertificateService interface {
	GetOrCreate(name string) (*Certificate, error)
}

// Service represents a Certificate service.
type Service struct {
	certRepo    CertificateRepository
	certFactory CertificateFactory
}

// NewCertificateService create new Certificate service.
func NewCertificateService(certRepo CertificateRepository, certFactory CertificateFactory) *Service {
	return &Service{
		certRepo:    certRepo,
		certFactory: certFactory,
	}
}

// GetOrCreate retrives or create a certificat for the given name.
func (s *Service) GetOrCreate(name string) (*Certificate, error) {
	var cert *Certificate
	cert, err := s.certRepo.Get(name)
	if errors.Is(err, ErrCertificateNotFound) {
		// Create new Certificate
		cert, err = s.certFactory.Create(name)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.GetOrCreate")
		}
		// Store Certificate
		err := s.certRepo.Store(cert)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.GetOrCreate")
		}

	}
	return cert, nil
}

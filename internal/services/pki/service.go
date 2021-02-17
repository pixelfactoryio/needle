package pki

import (
	"github.com/pkg/errors"
)

// ErrCertificateNotFound unable to find certificate
var ErrCertificateNotFound = errors.New("Certificate Not Found")

// CertificateService interface
type CertificateService interface {
	GetOrCreate(name string) (*Certificate, error)
}

type Service struct {
	certRepo    CertificateRepository
	certFactory CertificateFactory
}

// NewCertificateService create new CertificateService
func NewCertificateService(certRepo CertificateRepository, certFactory CertificateFactory) *Service {
	return &Service{
		certRepo:    certRepo,
		certFactory: certFactory,
	}
}

func (s *Service) GetOrCreate(name string) (*Certificate, error) {
	var cert *Certificate
	cert, err := s.certRepo.Get(name)
	if errors.Cause(err) == ErrCertificateNotFound {
		// Create new Certificate
		cert, err = s.certFactory.Create(name)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.FindOrCreate")
		}
		// Store Certificate
		err := s.certRepo.Store(cert)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.FindOrCreate")
		}

	}
	return cert, nil
}

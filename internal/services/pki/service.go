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

type certificateService struct {
	certRepo    CertificateRepository
	certFactory CertificateFactory
}

// NewCertificateService create new CertificateService
func NewCertificateService(certRepo CertificateRepository, certFactory CertificateFactory) CertificateService {
	return &certificateService{
		certRepo:    certRepo,
		certFactory: certFactory,
	}
}

func (cs *certificateService) GetOrCreate(name string) (*Certificate, error) {
	var cert *Certificate
	cert, err := cs.certRepo.Get(name)
	if errors.Cause(err) == ErrCertificateNotFound {
		// Create new Certificate
		cert, err = cs.certFactory.Create(name)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.FindOrCreate")
		}
		// Store Certificate
		err := cs.certRepo.Store(cert)
		if err != nil {
			return nil, errors.Wrap(err, "pki.CertificateService.FindOrCreate")
		}

	}
	return cert, nil
}

package pki

// CertificateRepository interface.
type CertificateRepository interface {
	Get(name string) (*Certificate, error)
	Store(certificate *Certificate) error
}

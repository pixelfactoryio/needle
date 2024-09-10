package pki

// InternalCert represents a certificate.
type InternalCert struct {
	Name      string `json:"name" storm:"id"`
	CertPEM   []byte `json:"cert_pem"`
	KeyPEM    []byte `json:"key_pem"`
	CreatedAt int64  `json:"created_at"`
}

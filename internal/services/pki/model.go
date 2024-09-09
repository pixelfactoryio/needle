package pki

// Certificate represents a certificate.
type Certificate struct {
	Name      string `json:"name" storm:"id"`
	CertPEM   []byte `json:"cert_pem"`
	KeyPEM    []byte `json:"key_pem"`
	CreatedAt int64  `json:"created_at"`
}

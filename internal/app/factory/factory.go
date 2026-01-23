package factory

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"go.pixelfactory.io/needle/internal/app/pki"
)

var pemEncode = pem.Encode

// Factory represents the certificate factory.
type Factory struct {
	rootCA tls.Certificate
}

// New create certificateFactory.
func New(rootCA tls.Certificate) *Factory {
	return &Factory{rootCA: rootCA}
}

// Create creates a certificate.
func (f *Factory) Create(name string) (*pki.InternalCert, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	// Try to parse name as IP.
	ipAddresses := []net.IP{net.ParseIP("0.0.0.0"), net.ParseIP("127.0.0.1")}
	if ip := net.ParseIP(name); ip != nil {
		ipAddresses = append(ipAddresses, ip)
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(2, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		DNSNames:    []string{"localhost", name},
		IPAddresses: ipAddresses,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	ca, err := x509.ParseCertificate(f.rootCA.Certificate[0])
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, f.rootCA.PrivateKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	if encodeErr := pemEncode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); encodeErr != nil {
		return nil, encodeErr
	}

	certPrivKeyPEM := new(bytes.Buffer)
	if encodeErr := pemEncode(
		certPrivKeyPEM,
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey)}); encodeErr != nil {
		return nil, encodeErr
	}

	return &pki.InternalCert{
		Name:    name,
		CertPEM: certPEM.Bytes(),
		KeyPEM:  certPrivKeyPEM.Bytes(),
	}, nil
}

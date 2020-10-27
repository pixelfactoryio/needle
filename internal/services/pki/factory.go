package pki

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
)

// CertificateFactory interface
type CertificateFactory interface {
	Create(name string) (*Certificate, error)
}

type factory struct {
	rootCA tls.Certificate
}

// NewFactory create certificateFactory
func NewFactory(rootCA tls.Certificate) CertificateFactory {
	return &factory{rootCA: rootCA}
}

// Create creates certificate
func (f *factory) Create(name string) (*Certificate, error) {
	serialNumber, err := getSerialNumber()
	if err != nil {
		return nil, err
	}

	// Try to parse name as IP
	IPAddresses := []net.IP{net.ParseIP("0.0.0.0"), net.ParseIP("127.0.0.1")}
	if ip := net.ParseIP(name); ip != nil {
		IPAddresses = append(IPAddresses, ip)
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
		IPAddresses: IPAddresses,
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
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return &Certificate{
		Name:    name,
		CertPEM: certPEM.Bytes(),
		KeyPEM:  certPrivKeyPEM.Bytes(),
	}, nil
}

func getSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	return serialNumber, nil
}

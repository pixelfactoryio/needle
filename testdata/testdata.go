package testdata

import (
	"crypto/tls"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"go.pixelfactory.io/needle/internal/app/pki"
)

// Dir returns testdata path
func Dir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func Setup(t *testing.T) (tls.Certificate, *pki.InternalCert) {
	rootCA, err := tls.LoadX509KeyPair(
		Dir()+"/certs/root-ca.crt",
		Dir()+"/certs/root-ca.key",
	)
	if err != nil {
		t.Error("Unable to get rootCA", err)
	}

	certPEM, err := ioutil.ReadFile(Dir() + "/certs/test.needle.local.crt")
	if err != nil {
		t.Error("Unable to get certPEM", err)
	}

	keyPEM, err := ioutil.ReadFile(Dir() + "/certs/test.needle.local.key")
	if err != nil {
		t.Error("Unable to get keyPEM", err)
	}

	testCert := pki.InternalCert{
		Name:    "test.needle.local",
		CertPEM: certPEM,
		KeyPEM:  keyPEM,
	}
	return rootCA, &testCert
}

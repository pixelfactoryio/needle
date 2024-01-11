// Package boltdb provides BoltDB repository.
package boltdb

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
	"go.pixelfactory.io/needle/internal/services/pki"
)

type boltRepository struct {
	client *storm.DB
}

// NewBoltRepository create BoltDB backed CertificateRepository
func NewBoltRepository(client *storm.DB) pki.CertificateRepository {
	repo := &boltRepository{
		client: client,
	}
	return repo
}

// Get certificate in data/cache.db
func (br *boltRepository) Get(name string) (*pki.Certificate, error) {
	var cert pki.Certificate
	err := br.client.One("Name", name, &cert)
	if err == storm.ErrNotFound {
		return nil, errors.Wrap(pki.ErrCertificateNotFound, "repository.BoltRepository.Find")
	}
	if err != nil {
		return nil, errors.Wrap(err, "repository.BoltRepository.Find")
	}
	return &cert, nil
}

// Store certificate in data/cache.db
func (br *boltRepository) Store(certificate *pki.Certificate) error {
	err := br.client.Save(certificate)
	if err != nil {
		return errors.Wrap(err, "repository.BoltRepository.store")
	}
	return nil
}

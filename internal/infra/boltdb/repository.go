// Package boltdb provides BoltDB repository.
package boltdb

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
	"go.pixelfactory.io/needle/internal/app/pki"
)

type boltRepository struct {
	client *storm.DB
}

// New creates a new BoltDB repository.
func New(client *storm.DB) pki.Repository {
	return &boltRepository{
		client: client,
	}
}

// Get certificate in data/cache.db.
func (br *boltRepository) Get(name string) (*pki.InternalCert, error) {
	var cert pki.InternalCert
	err := br.client.One("Name", name, &cert)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, errors.Wrap(pki.ErrCertificateNotFound, "repository.BoltRepository.Find")
	}
	if err != nil {
		return nil, errors.Wrap(err, "repository.BoltRepository.Find")
	}
	return &cert, nil
}

// Store certificate in data/cache.db.
func (br *boltRepository) Store(certificate *pki.InternalCert) error {
	err := br.client.Save(certificate)
	if err != nil {
		return errors.Wrap(err, "repository.BoltRepository.store")
	}
	return nil
}

package store

import (
	"bitbucket.org/heindl/provision/mgoeco"
	"time"
	"github.com/jonboulle/clockwork"
	"cloud.google.com/go/datastore"
	"bitbucket.org/heindl/provision/dseco"
)

type TaxaStore interface {
	ReadTaxa() (Taxa, error)
	ReadSpecies() (Taxa, error)
	NewIterator() *datastore.Iterator
	ReadTaxaFromCanonicalNames(...CanonicalName) (Taxa, error)
	SetTaxa(Taxa) error
	SetPhotos(Photos) error
	SetSchema(Schema) error
	UpdateSchemaLastFetched(Schema) error
	GetOccurrenceSchema(taxonKey *datastore.Key) (Schema, error)
	Close()
}

const SpeciesColl = mgoeco.CollectionName("species")


var _ TaxaStore = &store{}

func NewTestSpeciesStore() TaxaStore {

	client, err := dseco.NewMockDatastore()
	if err != nil {
		return nil
	}

	return TaxaStore(&store{clockwork.NewFakeClockAt(time.Now()), client})
}

func NewSpeciesStore() (TaxaStore, error) {

	client, err := dseco.NewLiveDatastore()
	if err != nil {
		return nil, err
	}

	return TaxaStore(&store{clockwork.NewRealClock(), client}), nil
}

type store struct {
	Clock clockwork.Clock
	DatastoreClient *datastore.Client
}

func (Ω *store) Close() {
	return
}


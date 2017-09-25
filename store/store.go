package store

import (
	"time"
	"github.com/jonboulle/clockwork"
	"cloud.google.com/go/datastore"
	"bitbucket.org/heindl/provision/dseco"
)

type TaxaStore interface {
	ReadTaxa() (Taxa, error)
	ReadSpecies() (Taxa, error)
	NewOccurrenceSchemeIterator(*datastore.Key) *datastore.Iterator
	ReadTaxaFromCanonicalNames(TaxonRank, ...CanonicalName) (Taxa, error)
	GetTaxon(*datastore.Key) (*Taxon, error)
	SetTaxa(Taxa) error
	SetPhotos(Photos) error
	SetSchema(Schema) error
	UpdateSchemaLastFetched(Schema) error
	GetOccurrenceSchema(*datastore.Key) (Schema, error)
	SetOccurrences(Occurrences) error
	GetOccurrenceIterator(taxonKey *datastore.Key) *datastore.Iterator
	GetOccurrences(taxonKey *datastore.Key) (Occurrences, error)
	Close() error
}


var _ TaxaStore = &store{}

func NewTestTaxaStore() TaxaStore {

	client, err := dseco.NewMockDatastore()
	if err != nil {
		return nil
	}

	return TaxaStore(&store{clockwork.NewFakeClockAt(time.Now()), client})
}

func NewTaxaStore() (TaxaStore, error) {

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

func (Ω *store) Close() error {
	return Ω.DatastoreClient.Close()
}


package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/jonboulle/clockwork"
	"time"
)

type TaxaStore interface {
	ReadTaxa(context.Context) (Taxa, error)
	ReadSpecies(context.Context) (Taxa, error)
	ReadTaxaFromCanonicalNames(context.Context, TaxonRank, ...CanonicalName) (Taxa, error)
	ReadTaxon(context.Context, INaturalistTaxonID) (*Taxon, error)
	CreateTaxonIfNotExists(context.Context, Taxon) error
	SetTaxonPhoto(context.Context, INaturalistTaxonID, string) error
	SetPhoto(context.Context, Photo) error
	UpsertDataSource(context.Context, DataSource) error
	GetSourceLastCreated(cxt context.Context, kind DataSourceKind, srcID DataSourceType) (*time.Time, error)
	UpdateDataSourceLastFetched(context.Context, DataSource) error
	GetOccurrenceDataSources(context.Context, INaturalistTaxonID) (DataSources, error)
	SetPrediction(cxt context.Context, p Prediction) error
	Close() error
}

var _ TaxaStore = &store{}

func NewTestTaxaStore() TaxaStore {

	s := store{
		Clock: clockwork.NewFakeClockAt(time.Now()),
	}

	var err error

	s.FirestoreClient, err = NewMockFirestore()
	if err != nil {
		panic(err)
	}

	return TaxaStore(&s)
}

func NewTaxaStore() (TaxaStore, error) {

	s := store{
		Clock: clockwork.NewRealClock(),
	}

	var err error

	s.FirestoreClient, err = NewLiveFirestore()
	if err != nil {
		return nil, err
	}

	return TaxaStore(&s), nil
}

type store struct {
	Clock                clockwork.Clock
	FirestoreClient      *firestore.Client
}

func (Ω *store) Close() error {
	return Ω.FirestoreClient.Close()
}

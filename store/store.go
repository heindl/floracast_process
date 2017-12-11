package store

import (
	"time"
	"github.com/jonboulle/clockwork"
	"cloud.google.com/go/firestore"
	"context"
)

type TaxaStore interface {
	ReadTaxa(context.Context) (Taxa, error)
	ReadSpecies(context.Context) (Taxa, error)
	ReadTaxaFromCanonicalNames(context.Context, TaxonRank, ...CanonicalName) (Taxa, error)
	GetTaxon(context.Context, TaxonID) (*Taxon, error)
	CreateTaxonIfNotExists(context.Context, Taxon) error
	SetTaxonPhoto(context.Context, TaxonID, string) error
	IncrementTaxonEcoRegion(cxt context.Context, taxonID TaxonID, ecoRegionKey string) error
	SetPhoto(context.Context, Photo) error
	UpsertDataSource(context.Context, DataSource) error
	UpdateDataSourceLastFetched(context.Context, DataSource) error
	GetOccurrenceDataSources(context.Context, TaxonID) (DataSources, error)
	UpsertOccurrence(context.Context, Occurrence) (isNewOccurrence bool, err error)
	GetOccurrences(context.Context, TaxonID) (Occurrences, error)
	ReadProtectedArea(cxt context.Context, lat, lng float64) (*ProtectedArea, error)
	SetProtectedArea(context.Context, ProtectedArea) error
	SetPrediction(cxt context.Context, p Prediction) error
	Close() error
}

var _ TaxaStore = &store{}

func NewTestTaxaStore() TaxaStore {

	client, err := NewMockFirestore()
	if err != nil {
		return nil
	}

	return TaxaStore(&store{clockwork.NewFakeClockAt(time.Now()), client})
}

func NewTaxaStore() (TaxaStore, error) {

	client, err := NewLiveFirestore()
	if err != nil {
		return nil, err
	}

	return TaxaStore(&store{clockwork.NewRealClock(), client}), nil
}

type store struct {
	Clock          clockwork.Clock
	FirestoreClient *firestore.Client
}

func (Ω *store) Close() error {
	return Ω.FirestoreClient.Close()
}


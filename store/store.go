package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

type FloraStore interface {
	FirestoreCollection(FirestoreCollection) *firestore.CollectionRef
	FirestoreBatch() *firestore.WriteBatch
	FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error
	AlgoliaIndex(indexFunc AlgoliaIndexFunc) (AlgoliaIndex, error)
	Close() error
}

func NewFloraStore(ctx context.Context) (FloraStore, error) {
	//s := store{
	//	Clock: clockwork.NewRealClock(),
	//}
	firestoreClient, err := NewLiveFirestore(ctx)
	if err != nil {
		return nil, err
	}

	algoliaClient, err := NewLiveAlgoliaClient(ctx)
	if err != nil {
		return nil, err
	}

	return &store{
		firestoreClient: firestoreClient,
		algoliaClient: algoliaClient,
	}, nil
}

type store struct {
	firestoreClient      *firestore.Client
	algoliaClient algoliasearch.Client
}

func (Ω *store) FirestoreCollection(æ FirestoreCollection) *firestore.CollectionRef {
	return Ω.firestoreClient.Collection(string(æ))
}


func (Ω *store) FirestoreBatch() *firestore.WriteBatch {
	return Ω.firestoreClient.Batch()
}

type FirestoreTransactionFunc func(context.Context, *firestore.Transaction) error

func (Ω *store) FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error {
	return Ω.firestoreClient.RunTransaction(ctx, fn)
}


func (Ω *store) AlgoliaIndex(æ AlgoliaIndexFunc) (AlgoliaIndex, error) {
	return æ(Ω.algoliaClient)
}

func (Ω *store) Close() error {
	return Ω.firestoreClient.Close()
}

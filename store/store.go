package store

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

type FloraStore interface {
	FirestoreCollection(FirestoreCollection) (*firestore.CollectionRef, error)
	FirestoreBatch() *firestore.WriteBatch
	FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error
	AlgoliaIndex(indexFunc AlgoliaIndexFunc) (AlgoliaIndex, error)
	CloudStorageBucket() (*storage.BucketHandle, error)
	Close() error
}

type TestFloraStore interface {
	FloraStore
	CountTestCollection(context.Context, *firestore.CollectionRef) (int, error)
	ClearTestCollection(context.Context, *firestore.CollectionRef) error
}

const (
	GCSTestBucket = "floracast-datamining-test"
	GCSLiveBucket = "floracast-datamining"
)

func NewTestFloraStore(ctx context.Context) (TestFloraStore, error) {
	//s := store{
	//	Clock: clockwork.NewRealClock(),
	//}
	firestoreClient, err := newLiveFirestore(ctx)
	if err != nil {
		return nil, err
	}

	algoliaClient, err := newLiveAlgoliaClient(ctx)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Google Cloud Storage client.")
	}

	gcsBucketHandle := client.Bucket(GCSTestBucket)
	attrs, err := gcsBucketHandle.Attrs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve GCS Bucket [%s] Attrs", GCSTestBucket)
	}
	if attrs.Created.IsZero() {
		return nil, errors.Newf("Bucket doesn't exist [%s]", GCSTestBucket)
	}

	return &store{
		gcsBucketHandle: gcsBucketHandle,
		isTest:          true,
		firestoreClient: firestoreClient,
		algoliaClient:   algoliaClient,
	}, nil
}

func NewFloraStore(ctx context.Context) (FloraStore, error) {
	//s := store{
	//	Clock: clockwork.NewRealClock(),
	//}
	firestoreClient, err := newLiveFirestore(ctx)
	if err != nil {
		return nil, err
	}

	algoliaClient, err := newLiveAlgoliaClient(ctx)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Google Cloud Storage client.")
	}

	gcsBucketHandle := client.Bucket(GCSLiveBucket)
	attrs, err := gcsBucketHandle.Attrs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve GCS Bucket [%s] Attrs", GCSLiveBucket)
	}
	if attrs.Created.IsZero() {
		return nil, errors.Newf("Bucket doesn't exist [%s]", GCSLiveBucket)
	}

	return &store{
		gcsBucketHandle: gcsBucketHandle,
		firestoreClient: firestoreClient,
		algoliaClient:   algoliaClient,
	}, nil
}

type store struct {
	isTest          bool
	firestoreClient *firestore.Client
	algoliaClient   algoliasearch.Client
	gcsBucketHandle *storage.BucketHandle
}

func (Ω *store) FirestoreCollection(æ FirestoreCollection) (*firestore.CollectionRef, error) {
	name := string(æ)
	if Ω.isTest {
		name = "Test" + name
	}
	return Ω.firestoreClient.Collection(name), nil
}

func (Ω *store) FirestoreBatch() *firestore.WriteBatch {
	return Ω.firestoreClient.Batch()
}

type FirestoreTransactionFunc func(context.Context, *firestore.Transaction) error

func (Ω *store) FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error {
	return Ω.firestoreClient.RunTransaction(ctx, fn)
}

func (Ω *store) AlgoliaIndex(æ AlgoliaIndexFunc) (AlgoliaIndex, error) {
	return æ(Ω.algoliaClient, Ω.isTest)
}

func (Ω *store) CloudStorageBucket() (*storage.BucketHandle, error) {
	return Ω.gcsBucketHandle, nil
}

func (Ω *store) CountTestCollection(ctx context.Context, col *firestore.CollectionRef) (int, error) {
	if !strings.Contains(strings.ToLower(col.ID), "test") {
		return 0, errors.Newf("Collection [%s] should include 'Test' in name", col.ID)
	}
	snaps, err := col.Documents(ctx).GetAll()
	if err != nil {
		return 0, errors.Newf("Could not count test collection [%s] documents", col.ID)
	}
	return len(snaps), nil
}

func (Ω *store) ClearTestCollection(ctx context.Context, col *firestore.CollectionRef) error {
	if !strings.Contains(strings.ToLower(col.ID), "test") {
		return errors.Newf("Collection [%s] should include 'Test' in name", col.ID)
	}
	snaps, err := col.Documents(ctx).GetAll()
	if err != nil {
		return errors.Wrapf(err, "Could get test collection [%s] documents", col.ID)
	}
	for _, snap := range snaps {
		if _, err := snap.Ref.Delete(ctx); err != nil {
			return errors.Wrapf(err, "Could not delete test collection [%s] record [%]", col.ID, snap.Ref.ID)
		}
	}
	return nil
}

func (Ω *store) Close() error {
	return Ω.firestoreClient.Close()
}

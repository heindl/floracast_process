package store

import (
	"bitbucket.org/heindl/process/utils"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"strings"
)

// FloraStore is the standard interface for updating the data store.
type FloraStore interface {
	FirestoreCollection(FireStoreCollection) (*firestore.CollectionRef, error)
	FirestoreBatch() *firestore.WriteBatch
	FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error
	AlgoliaIndex(indexFunc AlgoliaIndexFunc) (AlgoliaIndex, error)
	//CloudStorageBucket() (*storage.BucketHandle, error)
	CloudStorageObjects(ctx context.Context, pathname string, requiredSuffixes ...string) ([]*storage.ObjectHandle, error)
	Close() error
}

// TestFloraStore adds additional methods for testing.
type TestFloraStore interface {
	FloraStore
	CountTestCollection(context.Context, *firestore.CollectionRef) (int, error)
	ClearTestCollection(context.Context, *firestore.CollectionRef) error
}

const (
	gcsTestBucket = "floracast-datamining-test"
	gcsLiveBucket = "floracast-datamining"
)

// NewTestFloraStore returns an interface that guarantees writes to test collections.
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

	gcsBucketHandle := client.Bucket(gcsTestBucket)
	attrs, err := gcsBucketHandle.Attrs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve GCS Bucket [%s] Attrs", gcsTestBucket)
	}
	if attrs.Created.IsZero() {
		return nil, errors.Newf("Bucket doesn't exist [%s]", gcsTestBucket)
	}

	return &store{
		gcsBucketHandle: gcsBucketHandle,
		isTest:          true,
		firestoreClient: firestoreClient,
		algoliaClient:   algoliaClient,
	}, nil
}

// NewFloraStore initializes a database connections.
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

	gcsBucketHandle := client.Bucket(gcsLiveBucket)
	attrs, err := gcsBucketHandle.Attrs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not retrieve GCS Bucket [%s] Attrs", gcsLiveBucket)
	}
	if attrs.Created.IsZero() {
		return nil, errors.Newf("Bucket doesn't exist [%s]", gcsLiveBucket)
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

func (Ω *store) FirestoreCollection(æ FireStoreCollection) (*firestore.CollectionRef, error) {
	name := string(æ)
	if Ω.isTest {
		name = "Test" + name
	}
	return Ω.firestoreClient.Collection(name), nil
}

func (Ω *store) FirestoreBatch() *firestore.WriteBatch {
	return Ω.firestoreClient.Batch()
}

// FirestoreTransactionFunc is a callback for a transaction
type FirestoreTransactionFunc func(context.Context, *firestore.Transaction) error

func (Ω *store) FirestoreTransaction(ctx context.Context, fn FirestoreTransactionFunc) error {
	return Ω.firestoreClient.RunTransaction(ctx, fn)
}

func (Ω *store) AlgoliaIndex(æ AlgoliaIndexFunc) (AlgoliaIndex, error) {
	return æ(Ω.algoliaClient, Ω.isTest)
}

// ".tfrecord.gz"
func (Ω *store) CloudStorageObjects(ctx context.Context, pathname string, requiredSuffixes ...string) ([]*storage.ObjectHandle, error) {

	pathname = strings.TrimSpace(pathname)
	if strings.HasPrefix(pathname, "gs://") {
		pathname = strings.TrimPrefix(pathname, "gs://")
		l := strings.Split(pathname, "/")
		pathname = strings.Join(l[1:], "/")
	}

	if len(requiredSuffixes) > 0 && utils.HasSuffix(pathname, requiredSuffixes...) {
		return []*storage.ObjectHandle{Ω.gcsBucketHandle.Object(pathname)}, nil
	}

	iter := Ω.gcsBucketHandle.Objects(ctx, &storage.Query{
		Prefix: pathname,
	})

	gcsObjects := []*storage.ObjectHandle{}

	for {
		o, err := iter.Next()
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "Could not iterate over object names")
		}
		if len(requiredSuffixes) == 0 || utils.HasSuffix(o.Name, requiredSuffixes...) {
			gcsObjects = append(gcsObjects, Ω.gcsBucketHandle.Object(o.Name))
		}
	}

	return gcsObjects, nil
}

//func (Ω *store) CloudStorageBucket() (*storage.BucketHandle, error) {
//	return Ω.gcsBucketHandle, nil
//}

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

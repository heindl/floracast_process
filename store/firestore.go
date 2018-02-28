package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/option"
	"os"
	"time"
)

type FirestoreCollection string

const (
	CollectionOccurrences    = FirestoreCollection("Occurrences")
	CollectionRandom         = FirestoreCollection("Random")
	CollectionTaxa           = FirestoreCollection("Taxa")
	CollectionNameUsages     = FirestoreCollection("NameUsages")
	CollectionPhotos         = FirestoreCollection("Photos")
	CollectionProtectedAreas = FirestoreCollection("ProtectedAreas")
	CollectionPredictions    = FirestoreCollection("Predictions")
)

func NewFirestoreLimiter() <-chan time.Time {
	// Maximum writes per second per database (at beta): 2,500 (up to 2.5 MiB per second)
	t := time.NewTicker(time.Second / 1000)
	return t.C
}

// gcloud beta emulators datastore start --project=floracast-20c01 --store-on-disk=false

func newLiveFirestore(ctx context.Context) (*firestore.Client, error) {

	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("FIRESTORE_PROJECT_ID invalid")
	}

	var opts []option.ClientOption
	//if key := os.Getenv("FLORACAST_GCLOUD_API_KEY"); key != "" {
	//	opts = append(opts, option.WithAPIKey(key))
	//}

	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "could not get client")
	}

	return client, nil
}

// In the short term just make this a test project.
func NewMockFirestore(ctx context.Context) (*firestore.Client, error) {
	return nil, nil

	//projectID := os.Getenv("FLORACAST_GCLOUD_PROJECT_ID")
	//if projectID == "" {
	//	return nil, errors.New("FLORACAST_GCLOUD_PROJECT_ID invalid")
	//}
	//
	//client, err := datastore.NewClient(context.Background(), projectID, option.WithEndpoint("localhost:8081"))
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not get client")
	//}
	//return client, nil
}

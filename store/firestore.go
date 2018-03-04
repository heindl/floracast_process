package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/option"
	"os"
	"time"
)

// FireStoreCollection is the collection name
type FireStoreCollection string

const (
	// CollectionOccurrences is a constant for the FireStore Collection name
	CollectionOccurrences = FireStoreCollection("Occurrences")
	// CollectionRandom is a constant for the FireStore Collection name
	CollectionRandom = FireStoreCollection("Random")
	// CollectionTaxa is a constant for the FireStore Collection name
	CollectionTaxa = FireStoreCollection("Taxa")
	// CollectionNameUsages is a constant for the FireStore Collection name
	CollectionNameUsages = FireStoreCollection("NameUsages")
	// CollectionProtectedAreas is a constant for the FireStore Collection name
	CollectionProtectedAreas = FireStoreCollection("ProtectedAreas")
	// CollectionPredictions is a constant for the FireStore Collection name
	CollectionPredictions = FireStoreCollection("Predictions")
)

// NewFirestoreLimiter is a reminder for the necessary FireStore rate limit.
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

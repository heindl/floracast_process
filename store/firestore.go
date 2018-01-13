package store

import (
	"cloud.google.com/go/firestore"
	"github.com/saleswise/errors/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"
)

// gcloud beta emulators datastore start --project=floracast-20c01 --store-on-disk=false

func NewLiveFirestore() (*firestore.Client, error) {

	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("FIRESTORE_PROJECT_ID invalid")
	}

	var opts []option.ClientOption
	//if key := os.Getenv("FLORACAST_GCLOUD_API_KEY"); key != "" {
	//	opts = append(opts, option.WithAPIKey(key))
	//}

	client, err := firestore.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "could not get client")
	}

	return client, nil
}

// In the short term just make this a test project.
func NewMockFirestore() (*firestore.Client, error) {
	return NewLiveFirestore()

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

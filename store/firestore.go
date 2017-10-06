package store

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"os"
	"github.com/saleswise/errors/errors"
	"google.golang.org/api/option"
)

// gcloud beta emulators datastore start --project=floracast-20c01 --store-on-disk=false

func NewLiveDatastore() (*firestore.Client, error) {

	projectID := os.Getenv("FLORACAST_GCLOUD_PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("FLORACAST_GCLOUD_PROJECT_ID invalid")
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
func NewMockDatastore() (*firestore.Client, error) {
	return NewLiveDatastore()


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
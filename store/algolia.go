package store

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"os"
	"context"
)

const envAPIKey = "envAPIKey"
const envApplicationID = "envApplicationID"


type AlgoliaIndexFunc func(client algoliasearch.Client) (algoliasearch.Index, error)

type AlgoliaIndex interface{
	AddObjects([]algoliasearch.Object) (algoliasearch.BatchRes, error)
	DeleteBy(params algoliasearch.Map) (res algoliasearch.DeleteTaskRes, err error)
}

func NewLiveAlgoliaClient(ctx context.Context) (algoliasearch.Client, error) {

	if os.Getenv(envAPIKey) == "" || os.Getenv(envApplicationID) == "" {
		return nil, errors.Newf("%s and %s environment variables required for an Algolia Index", envAPIKey, envApplicationID)
	}

	client := algoliasearch.NewClient(os.Getenv(envApplicationID), os.Getenv(envAPIKey))

	return client, nil
}

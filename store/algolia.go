package store

import (
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"os"
)

const algoliaEnvAPIKey = "FLORACAST_ALGOLIA_API_KEY"
const algoliaEnvApplicationID = "FLORACAST_ALGOLIA_APPLICATION_ID"

type AlgoliaIndexFunc func(client algoliasearch.Client, isTest bool) (algoliasearch.Index, error)

type AlgoliaIndex interface {
	AddObjects([]algoliasearch.Object) (algoliasearch.BatchRes, error)
	BrowseAll(params algoliasearch.Map) (it algoliasearch.IndexIterator, err error)
	DeleteBy(params algoliasearch.Map) (res algoliasearch.DeleteTaskRes, err error)
}

func newLiveAlgoliaClient(ctx context.Context) (algoliasearch.Client, error) {

	if os.Getenv(algoliaEnvAPIKey) == "" || os.Getenv(algoliaEnvApplicationID) == "" {
		return nil, errors.Newf("%s and %s environment variables required for an Algolia Index", algoliaEnvAPIKey, algoliaEnvApplicationID)
	}

	client := algoliasearch.NewClient(os.Getenv(algoliaEnvApplicationID), os.Getenv(algoliaEnvAPIKey))

	return client, nil
}

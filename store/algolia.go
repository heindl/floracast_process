package store

import (
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"os"
)

const algoliaEnvAPIKey = "FLORACAST_ALGOLIA_API_KEY"
const algoliaEnvApplicationID = "FLORACAST_ALGOLIA_APPLICATION_ID"

// AlgoliaIndexFunc is a callback for retrieving an Algolia Index
type AlgoliaIndexFunc func(client algoliasearch.Client, isTest bool) (algoliasearch.Index, error)

// AlgoliaIndex is an interface for updating Algolia
type AlgoliaIndex interface {
	DeleteObjects(objectIDs []string) (algoliasearch.BatchRes, error)
	UpdateObject(object algoliasearch.Object) (res algoliasearch.UpdateObjectRes, err error)
	AddObjects([]algoliasearch.Object) (algoliasearch.BatchRes, error)
	BrowseAll(params algoliasearch.Map) (it algoliasearch.IndexIterator, err error)
	DeleteBy(params algoliasearch.Map) (res algoliasearch.DeleteTaskRes, err error)
}

type AlgoliaIndexName string

func (Ω *store) AlgoliaIndex(name AlgoliaIndexName, settings algoliasearch.Map) (AlgoliaIndex, error) {

	if Ω.isTest {
		name = AlgoliaIndexName(fmt.Sprintf("Test%s", name))
	}

	knownIndexes, err := Ω.algoliaClient.ListIndexes()
	if err != nil {
		return nil, errors.Wrap(err, "Could not access Algolia indexes")
	}

	for _, i := range knownIndexes {
		if i.Name == string(name) {
			// Ignore settings if the index already exists.
			return AlgoliaIndex(Ω.algoliaClient.InitIndex(i.Name)), nil
		}
	}

	index := Ω.algoliaClient.InitIndex(string(name))

	if settings != nil {
		if _, err := index.SetSettings(settings); err != nil {
			return nil, errors.Wrap(err, "Could not add settings to NameUsage Algolia index")
		}
	}

	return AlgoliaIndex(index), nil
}

func newLiveAlgoliaClient(ctx context.Context) (algoliasearch.Client, error) {

	if os.Getenv(algoliaEnvAPIKey) == "" || os.Getenv(algoliaEnvApplicationID) == "" {
		return nil, errors.Newf("%s and %s environment variables required for an Algolia Index", algoliaEnvAPIKey, algoliaEnvApplicationID)
	}

	client := algoliasearch.NewClient(os.Getenv(algoliaEnvApplicationID), os.Getenv(algoliaEnvAPIKey))

	return client, nil
}

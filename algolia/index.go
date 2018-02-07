package algolia

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"os"
	"github.com/saleswise/errors/errors"
	"fmt"
)


const indexNameUsage = "NameUsage"

type AlgoliaIndex interface{
	AddObjects([]algoliasearch.Object) (algoliasearch.BatchRes, error)
	DeleteBy(params algoliasearch.Map) (res algoliasearch.DeleteTaskRes, err error)
}

const envAPIKey = "envAPIKey"
const envApplicationID = "envApplicationID"

func NewAlgoliaNameUsageIndex() (AlgoliaIndex, error) {

	if os.Getenv(envAPIKey) == "" || os.Getenv(envApplicationID) == "" {
		return nil, errors.Newf("%s and %s environment variables required for an Algolia Index", envAPIKey, envApplicationID)
	}

	client := algoliasearch.NewClient(os.Getenv(envApplicationID), os.Getenv(envAPIKey))

	index := client.InitIndex(indexNameUsage)

	if _, err := index.SetSettings(algoliasearch.Map{
		"distinct": keyNameUsageID,
		"customRanking": []string{
			fmt.Sprintf("desc(%s)", keyReferenceCount),
		},
		"searchableAttributes": []string{
			string(keyCommonName),
			string(keyScientificName),
		},
	}); err != nil {
		return nil, errors.Wrap(err, "Could not add settings to NameUsage Algolia index")
	}

	return index, nil

}
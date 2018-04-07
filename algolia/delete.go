package algolia

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
)

func DeleteNameUsage(ctx context.Context, florastore store.FloraStore, nameUsageID nameusage.ID) error {

	for _, i := range []store.AlgoliaIndexName{PredictionIndex, OccurrenceIndex, NameUsageIndex} {
		index, err := florastore.AlgoliaIndex(i, nil)
		if err != nil {
			return err
		}
		if _, err := index.DeleteBy(algoliasearch.Map{"NameUsageID": nameUsageID}); err != nil {
			return errors.Wrapf(err, "Could not delete records from Algolia Index [%s] for NameUsageID [%s]", i, nameUsageID)
		}
	}

	return nil

}

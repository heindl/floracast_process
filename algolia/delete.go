package algolia

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
)

func DeleteNameUsage(ctx context.Context, florastore store.FloraStore, nameUsageID nameusage.ID) error {

	index, err := florastore.AlgoliaIndex(NameUsageIndex, nil)
	if err != nil {
		return err
	}
	if _, err := index.DeleteBy(
		algoliasearch.Map{
			"filters": fmt.Sprintf("NameUsageID:%s", nameUsageID),
		}); err != nil {
		return errors.Wrapf(err, "Could not delete records from Algolia Index for NameUsageID [%s]", nameUsageID)
	}

	return nil

}

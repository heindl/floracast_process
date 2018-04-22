package algolia

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
)

type NameUsageIndexRecord struct {
	StandardNameUsageRecord
	ScientificNames []string `json:""`
	CommonNames     []string `json:""`
}

const NameUsageIndex = store.AlgoliaIndexName("NameUsages")

var nameUsageIndexSettings = algoliasearch.Map{
	"searchableAttributes": []string{
		"CommonNames",
		"ScientificNames",
	},
}

// UploadNameUsageObjects creates searchable Algolia objects.
func IndexNameUsage(ctx context.Context, floraStore store.FloraStore, nameUsageID nameusage.ID) error {

	// Need both the NameUsage and MaterializedTaxon, which contains the thumbnail link.
	// This is necessary because the MaterializedTaxon selects and possibly resizes the image.
	record, err := generateNameUsage(ctx, floraStore, nameUsageID)
	if err != nil {
		return err
	}
	return record.upload()
}

func generateNameUsage(ctx context.Context, floraStore store.FloraStore, nameUsageID nameusage.ID) (*NameUsageIndexRecord, error) {
	record := NameUsageIndexRecord{
		StandardNameUsageRecord: StandardNameUsageRecord{
			NameUsageID: nameUsageID,
			floraStore:  floraStore,
		},
	}

	if err := record.hydrate(); err != nil {
		return nil, err
	}

	if err := record.fetchNameReferences(); err != nil {
		return nil, err
	}

	return &record, nil
}

func (Ω *NameUsageIndexRecord) fetchNameReferences() error {
	usage, err := nameusage.Fetch(context.Background(), Ω.floraStore, Ω.NameUsageID)
	if err != nil {
		return err
	}

	commonNames, err := usage.CommonNameReferenceLedger()
	if err != nil {
		return err
	}
	Ω.CommonNames = commonNames.Names()

	scientificNames, err := usage.ScientificNameReferenceLedger()
	if err != nil {
		return err
	}
	Ω.ScientificNames = scientificNames.Names()

	return nil
}

func (Ω *NameUsageIndexRecord) upload() error {

	index, err := Ω.floraStore.AlgoliaIndex(NameUsageIndex, nameUsageIndexSettings)
	if err != nil {
		return err
	}

	predictionObject, err := asAlgoliaObject(Ω)
	if err != nil {
		return err
	}

	if _, err := index.AddObjects([]algoliasearch.Object{predictionObject}); err != nil {
		return errors.Wrap(err, "Could not add NameUsageIndexRecord to Algolia")
	}

	return nil

}

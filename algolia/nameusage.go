package algolia

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	"fmt"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"strings"
	"bitbucket.org/heindl/processors/utils"
	"bitbucket.org/heindl/processors/store"
	"context"
)

func UploadNameUsageObjects(ctx context.Context, floracastStore store.FloraStore, usage *nameusage.NameUsage, deletedUsages nameusage.NameUsageIDs) error {

	index, err := floracastStore.AlgoliaIndex(nameUsageIndex)
	if err != nil {
		return err
	}

	if err := clearNameUsageObjects(index, deletedUsages...); err != nil {
		return err
	}
	
	objs, err := generateNameUsageObjects(ctx, usage)
	if err != nil {
		return err
	}

	for _, batch := range objs.batches(500) {
		if _, err := index.AddObjects(batch.asAlgoliaMapObjects()); err != nil {
			return errors.Wrap(err, "Could not add Angolia NameUsage objects")
		}
	}
	return nil

}

func clearNameUsageObjects(index store.AlgoliaIndex, nameUsageIDs ...nameusage.NameUsageID) error {

	for _, usageID := range nameUsageIDs {
		if usageID == "" {
			return errors.New("NameUsage ID required to delete Algolia objects")
		}
		if _, err := index.DeleteBy(algoliasearch.Map{
			string(KeyNameUsageID): usageID,
		}); err != nil {
			return errors.Wrap(err, "Could not delete Algolia NameUsageObjects")
		}
	}
	return nil
}


const IndexNameUsage = "NameUsage"

func nameUsageIndex(client algoliasearch.Client) (algoliasearch.Index, error) {

	index := client.InitIndex(IndexNameUsage)

	if _, err := index.SetSettings(algoliasearch.Map{
		"distinct": KeyNameUsageID,
		"customRanking": []string{
			fmt.Sprintf("desc(%s)", KeyReferenceCount),
		},
		"searchableAttributes": []string{
			string(KeyCommonName),
			string(KeyScientificName),
		},
	}); err != nil {
		return nil, errors.Wrap(err, "Could not add settings to NameUsage Algolia index")
	}

	return index, nil

}
const (
	KeyNameUsageID     = ObjectKey("NameUsageID")
	KeyScientificName  = ObjectKey("ScientificName")
	KeyCommonName      = ObjectKey("CommonName")
	KeyThumbnail       = ObjectKey("Thumbnail")
	KeyOccurrenceCount = ObjectKey("TotalOccurrenceCount")
	KeyReferenceCount  = ObjectKey("ReferenceCount")
)

func generateNameUsageObjects(ctx context.Context, usage *nameusage.NameUsage) (AlgoliaObjects, error) {

	if usage.TotalOccurrenceCount() == 0 {
		// Note that the algolia generation should only be called after occurrences fetched.
		// The occurrence count allows us to sort search results in Autocomplete.
		return nil, errors.New("Expected name usage provided to Algolia to have occurrences")
	}

	usageCommonName, err := usage.CommonName()
	if err != nil {
		return nil, err
	}
	usageCommonName = strings.Title(usageCommonName)

	usageOccurrenceCount := usage.TotalOccurrenceCount()

	// TODO: Generate thumbnail from image
	thumbnail := ""

	res := AlgoliaObjects{}

	for _, ref := range usage.ScientificNameReferenceLedger() {
		res = append(res, AlgoliaObject{
			KeyNameUsageID:     usage.ID(),
			KeyScientificName:  utils.CapitalizeString(ref.Name),
			KeyCommonName:      usageCommonName,
			KeyThumbnail:       thumbnail,
			KeyOccurrenceCount: usageOccurrenceCount,
			KeyReferenceCount:  ref.ReferenceCount,
		})
	}

	for _, ref := range usage.CommonNameReferenceLedger() {
		res = append(res, AlgoliaObject{
			KeyNameUsageID:     usage.ID(),
			KeyScientificName:  utils.CapitalizeString(usage.CanonicalName().ScientificName()),
			KeyCommonName:      strings.Title(ref.Name),
			KeyThumbnail:       thumbnail,
			KeyOccurrenceCount: usageOccurrenceCount,
			KeyReferenceCount:  ref.ReferenceCount,
		})
	}

	return res, nil

}
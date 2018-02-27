package algolia

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
)

// UploadNameUsageObjects creates searchable Algolia objects.
func UploadNameUsageObjects(ctx context.Context, florastore store.FloraStore, usage nameusage.NameUsage, deletedUsages ...nameusage.NameUsageID) error {

	if err := deleteNameUsageObjects(florastore, deletedUsages...); err != nil {
		return err
	}

	objs, err := generateNameUsageObjects(ctx, usage)
	if err != nil {
		return err
	}

	return uploadNameUsageObjects(florastore, objs)
}

func uploadNameUsageObjects(florastore store.FloraStore, objs objects) error {

	index, err := florastore.AlgoliaIndex(nameUsageIndex)
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

func deleteNameUsageObjects(florastore store.FloraStore, nameUsageIDs ...nameusage.NameUsageID) error {

	if len(nameUsageIDs) == 0 {
		return nil
	}

	index, err := florastore.AlgoliaIndex(nameUsageIndex)
	if err != nil {
		return err
	}
	if _, err := index.DeleteBy(nameUsageIDFilter(nameUsageIDs...)); err != nil {
		return errors.Wrap(err, "Could not delete Algolia NameUsage Objects")
	}

	return nil
}

func nameUsageIDFilter(nameUsageIDs ...nameusage.NameUsageID) algoliasearch.Map {
	if len(nameUsageIDs) == 0 {
		return nil
	}
	filterKeys := []string{}
	for _, id := range nameUsageIDs {
		filterKeys = append(filterKeys, fmt.Sprintf("%s:%s", keyNameUsageID, id))
	}
	return algoliasearch.Map{
		"filters": strings.Join(filterKeys, " OR "),
	}
}

const indexNameUsage = "NameUsage"
const indexTestNameUsage = "TestNameUsage"

func nameUsageIndex(client algoliasearch.Client, isTest bool) (algoliasearch.Index, error) {

	var index algoliasearch.Index
	if isTest {
		index = client.InitIndex(indexTestNameUsage)
	} else {
		index = client.InitIndex(indexNameUsage)
	}

	if _, err := index.SetSettings(algoliasearch.Map{
		"distinct":              true,
		"attributeForDistinct":  string(keyNameUsageID),
		"attributesForFaceting": []string{string(keyNameUsageID)},
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

const (
	keyNameUsageID     = objectKey("NameUsageID")
	keyScientificName  = objectKey("ScientificName")
	keyCommonName      = objectKey("CommonName")
	keyThumbnail       = objectKey("Thumbnail")
	keyOccurrenceCount = objectKey("TotalOccurrenceCount")
	keyReferenceCount  = objectKey("ReferenceCount")
)

func generateNameUsageObjects(ctx context.Context, usage nameusage.NameUsage) (objects, error) {

	usageOccurrenceCount, err := usage.Occurrences()
	if err != nil {
		return nil, err
	}

	if usageOccurrenceCount == 0 {
		// Note that the algolia generation should only be called after occurrences fetched.
		// The occurrence count allows us to sort search results in Autocomplete.
		return nil, errors.New("Expected name usage provided to Algolia to have occurrences")
	}

	usageCommonName, err := usage.CommonName()
	if err != nil {
		return nil, err
	}
	usageCommonName = strings.Title(usageCommonName)

	// TODO: Generate thumbnail from image
	thumbnail := ""

	res := objects{}

	sciNameRefLedger, err := usage.ScientificNameReferenceLedger()
	if err != nil {
		return nil, err
	}

	id, err := usage.ID()
	if err != nil {
		return nil, err
	}

	for _, ref := range sciNameRefLedger {
		res = append(res, object{
			keyNameUsageID:     id,
			keyScientificName:  utils.CapitalizeString(ref.Name),
			keyCommonName:      usageCommonName,
			keyThumbnail:       thumbnail,
			keyOccurrenceCount: usageOccurrenceCount,
			keyReferenceCount:  ref.ReferenceCount,
		})
	}

	commonNameRefLedger, err := usage.CommonNameReferenceLedger()
	if err != nil {
		return nil, err
	}

	for _, ref := range commonNameRefLedger {
		res = append(res, object{
			keyNameUsageID:     id,
			keyScientificName:  utils.CapitalizeString(usage.CanonicalName().ScientificName()),
			keyCommonName:      strings.Title(ref.Name),
			keyThumbnail:       thumbnail,
			keyOccurrenceCount: usageOccurrenceCount,
			keyReferenceCount:  ref.ReferenceCount,
		})
	}

	return res, nil

}

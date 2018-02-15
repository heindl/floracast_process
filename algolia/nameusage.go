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

func uploadNameUsageObjects(florastore store.FloraStore, objs AlgoliaObjects) error {

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
		filterKeys = append(filterKeys, fmt.Sprintf("%s:%s", KeyNameUsageID, id))
	}
	return algoliasearch.Map{
		"filters": strings.Join(filterKeys, " OR "),
	}
}

func countNameUsages(florastore store.FloraStore, nameUsageIDs ...nameusage.NameUsageID) (int, error) {
	index, err := florastore.AlgoliaIndex(nameUsageIndex)
	if err != nil {
		return 0, err
	}
	count := 0

	iter, err := index.BrowseAll(nameUsageIDFilter(nameUsageIDs...))
	if err != nil && err != algoliasearch.NoMoreHitsErr {
		return 0, errors.Wrap(err, "Could not browse Algolia NameUsage Index")
	}

	for {
		if _, err := iter.Next(); err != nil && err != algoliasearch.NoMoreHitsErr {
			return 0, err
		} else if err != nil && err == algoliasearch.NoMoreHitsErr {
			break
		}
		count += 1
	}
	return count, nil

}


const IndexNameUsage = "NameUsage"
const IndexTestNameUsage = "TestNameUsage"

func nameUsageIndex(client algoliasearch.Client, isTest bool) (algoliasearch.Index, error) {

	var index algoliasearch.Index
	if isTest {
		index = client.InitIndex(IndexTestNameUsage)
	} else {
		index = client.InitIndex(IndexNameUsage)
	}

	if _, err := index.SetSettings(algoliasearch.Map{
		"distinct": true,
		"attributeForDistinct": string(KeyNameUsageID),
		"attributesForFaceting": []string{string(KeyNameUsageID)},
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

func generateNameUsageObjects(ctx context.Context, usage nameusage.NameUsage) (AlgoliaObjects, error) {

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

	res := AlgoliaObjects{}

	sciNameRefLedger, err := usage.ScientificNameReferenceLedger()
	if err != nil {
		return nil, err
	}

	id, err := usage.ID()
	if err != nil {
		return nil, err
	}

	for _, ref := range sciNameRefLedger {
		res = append(res, AlgoliaObject{
			KeyNameUsageID:     id,
			KeyScientificName:  utils.CapitalizeString(ref.Name),
			KeyCommonName:      usageCommonName,
			KeyThumbnail:       thumbnail,
			KeyOccurrenceCount: usageOccurrenceCount,
			KeyReferenceCount:  ref.ReferenceCount,
		})
	}

	commonNameRefLedger, err := usage.CommonNameReferenceLedger()
	if err != nil {
		return nil, err
	}

	for _, ref := range commonNameRefLedger {
		res = append(res, AlgoliaObject{
			KeyNameUsageID:     id,
			KeyScientificName:  utils.CapitalizeString(usage.CanonicalName().ScientificName()),
			KeyCommonName:      strings.Title(ref.Name),
			KeyThumbnail:       thumbnail,
			KeyOccurrenceCount: usageOccurrenceCount,
			KeyReferenceCount:  ref.ReferenceCount,
		})
	}

	return res, nil

}
package materialtaxa

import (
"github.com/algolia/algoliasearch-client-go/algoliasearch"
"github.com/dropbox/godropbox/errors"
"math"
"strings"
"bitbucket.org/heindl/taxa/utils"
"os"
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

type AlgoliaNameObject map[nameObjectKey]interface{}
type AlgoliaNameUsageObjects []AlgoliaNameObject


func (Ω AlgoliaNameUsageObjects) hasCombination(scientificName, commonName string) bool {

	for _, o := range Ω {

		if o[keyScientificName] != scientificName {
			continue
		}
		if o[keyCommonName] != commonName {
			continue
		}
		return true
	}
	return false
}

func (Ω AlgoliaNameUsageObjects) asAlgoliaMapObjects() []algoliasearch.Object {
	res := []algoliasearch.Object{}
	for _, nameObject := range Ω {
		o := algoliasearch.Object{}
		for k, v := range nameObject {
			o[string(k)] = v
		}
		res = append(res, o)
	}
	return res
}

func (Ω AlgoliaNameUsageObjects) batches(maxBatchSize float64) []AlgoliaNameUsageObjects {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []AlgoliaNameUsageObjects{}
	for i := 0.0; i <= batchCount - 1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)
		if end > len(Ω) {
			end = len(Ω) - 1
		}
		res = append([]AlgoliaNameUsageObjects{}, Ω[start:end])
	}

	return res
}



func GenerateAlgoliaNameUsageObjects(usage *NameUsage) (AlgoliaNameUsageObjects, error) {

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

	res := AlgoliaNameUsageObjects{}

	for _, ref := range usage.ScientificNameReferenceLedger() {
		res = append(res, AlgoliaNameObject{
			keyNameUsageID:     usage.ID(),
			keyScientificName:  utils.CapitalizeString(ref.Name),
			keyCommonName:      usageCommonName,
			keyThumbnail:       thumbnail,
			keyOccurrenceCount: usageOccurrenceCount,
			keyReferenceCount:  ref.ReferenceCount,
		})
	}

	for _, ref := range usage.CommonNameReferenceLedger() {
		res = append(res, AlgoliaNameObject{
			keyNameUsageID:     usage.ID(),
			keyScientificName:  utils.CapitalizeString(usage.CanonicalName().ScientificName()),
			keyCommonName:      strings.Title(ref.Name),
			keyThumbnail:       thumbnail,
			keyOccurrenceCount: usageOccurrenceCount,
			keyReferenceCount:  ref.ReferenceCount,
		})
	}

	return res, nil

}

func (Ω AlgoliaNameUsageObjects) UploadObjects(index AlgoliaIndex) error {
	for _, batch := range Ω.batches(500) {
		if _, err := index.AddObjects(batch.asAlgoliaMapObjects()); err != nil {
			return errors.Wrap(err, "Could not add Angolia NameUsage objects")
		}
	}
	return nil
}

func DeleteAlgoliaNameUsageObjects(index AlgoliaIndex, nameUsageID string) error {
	if nameUsageID == "" {
		return errors.New("NameUsage ID required to delete Algolia objects")
	}
	if _, err := index.DeleteBy(algoliasearch.Map{
		string(keyNameUsageID): nameUsageID,
	}); err != nil {
		return errors.Wrap(err, "Could not delete Algolia NameUsageObjects")
	}
	return nil
}

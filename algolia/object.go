package algolia

import (
	"math"

	"bitbucket.org/heindl/process/store"
	"encoding/json"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
)

func asAlgoliaObject(i interface{}) (algoliasearch.Object, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal NameUsageIndexRecord to json")
	}
	res := algoliasearch.Object{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal NameUsageIndexRecord from json")
	}
	return res, nil

}

func upload(indexName store.AlgoliaIndexName) error {

	indexSettings := algoliasearch.Map{}

	if indexName == PredictionIndex {
		indexSettings = predictionIndexSettings
		Ω.Geolocations = Ω.predictionLocations
	}
	if indexName == OccurrenceIndex {
		indexSettings = occurrenceIndexSettings
		Ω.Geolocations = Ω.occurrenceLocations
	}

	index, err := Ω.floraStore.AlgoliaIndex(indexName, indexSettings)
	if err != nil {
		return err
	}

	predictionObject, err := Ω.asObject()
	if err != nil {
		return err
	}

	if _, err := index.AddObjects([]algoliasearch.Object{predictionObject}); err != nil {
		return errors.Wrap(err, "Could not add NameUsageIndexRecord to Algolia")
	}

	return nil

}

type objectKey string
type object map[objectKey]interface{}
type objects []object

func (Ω objects) asAlgoliaMapObjects() []algoliasearch.Object {
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

func (Ω objects) batches(maxBatchSize float64) []objects {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []objects{}
	for i := 0.0; i <= batchCount-1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)
		if end > len(Ω) {
			end = len(Ω)
		}
		o := Ω[start:end]
		res = append(res, o)
	}

	return res
}

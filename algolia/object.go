package algolia

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"math"
)

type ObjectKey string
type AlgoliaObject map[ObjectKey]interface{}
type AlgoliaObjects []AlgoliaObject

func (Ω AlgoliaObjects) hasCombination(scientificName, commonName string) bool {

	for _, o := range Ω {

		if o[KeyScientificName] != scientificName {
			continue
		}
		if o[KeyCommonName] != commonName {
			continue
		}
		return true
	}
	return false
}

func (Ω AlgoliaObjects) asAlgoliaMapObjects() []algoliasearch.Object {
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

func (Ω AlgoliaObjects) batches(maxBatchSize float64) []AlgoliaObjects {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []AlgoliaObjects{}
	for i := 0.0; i <= batchCount - 1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)
		if end > len(Ω) {
			end = len(Ω) - 1
		}
		res = append([]AlgoliaObjects{}, Ω[start:end])
	}

	return res
}

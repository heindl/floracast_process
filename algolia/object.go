package algolia

import (
	"math"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

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

package aggregate_usages

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"context"
	"fmt"
	"encoding/json"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of sources", t, func() {

		srcs, err := AggregateNameUsages(context.Background(), 47348, 56830, 48701)
		So(err, ShouldBeNil)
		////


		sufficient := name_usage.CanonicalNameUsages{}
		for _, src := range srcs {
			if src.OccurrenceCount() < 100 {
				continue
			}
			sufficient = append(sufficient, src)
		}


		marshalledSources, err := json.Marshal(sufficient)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(marshalledSources))

		//b, err := ioutil.ReadFile("/Users/m/Desktop/morchella.json")
		//So(err, ShouldBeNil)
		////
		//sources := CanonicalNameSources{}
		//So(json.Unmarshal(b, &sources), ShouldBeNil)
		//
		//for _, src := range res {
		//	fmt.Println(src.CanonicalName)
		//	fmt.Println(strings.Join(src.Synonyms, ", "))
		//	fmt.Println(src.SynonymFor)
		//	for k, v := range src.SourceTargetOccurrenceCount {
		//		fmt.Println(k, ":", len(v))
		//	}
		//	fmt.Println("-------------------")
		//	fmt.Println("-------------------")
		//}

	})
}

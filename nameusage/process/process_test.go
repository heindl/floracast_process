package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of sources", t, func() {


		// 47348, 56830, 48701
		srcs, err := AggregateNameUsages(context.Background(), 47348)
		So(err, ShouldBeNil)
		////
		fmt.Println(utils.JsonOrSpew(srcs))


		//sufficient := name_usage.Aggregate{}
		//for _, src := range srcs {
		//	if src.TotalOccurrenceCount() < 100 {
		//		continue
		//	}
		//	sufficient = append(sufficient, src)
		//}
		//
		//
		//marshalledSources, err := json.Marshal(sufficient)
		//if err != nil {
		//	panic(err)
		//}
		//
		//fmt.Println(string(marshalledSources))

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
		//	for k, v := range src.sourceTargetOccurrenceCount {
		//		fmt.Println(k, ":", len(v))
		//	}
		//	fmt.Println("-------------------")
		//	fmt.Println("-------------------")
		//}

	})
}

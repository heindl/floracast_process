package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"encoding/json"
	"io/ioutil"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of sources", t, func() {

		//srcs, err := fetchTaxaSources(context.Background(), inaturalist.ParseStringIDs("56830")...)
		//So(err, ShouldBeNil)
		////
		//marshalledSources, err := json.Marshal(srcs)
		//if err != nil {
		//	panic(err)
		//}
		//
		//fmt.Println(string(marshalledSources))

		b, err := ioutil.ReadFile("/Users/m/Desktop/morchella.json")
		So(err, ShouldBeNil)
		//
		sources := CanonicalNameSources{}
		So(json.Unmarshal(b, &sources), ShouldBeNil)
		//
		_ = sources.GenerateNameResults()
		//for _, src := range res {
		//	fmt.Println(src.CanonicalName)
		//	fmt.Println(strings.Join(src.Synonyms, ", "))
		//	fmt.Println(src.SynonymFor)
		//	for k, v := range src.SourceOccurrenceCount {
		//		fmt.Println(k, ":", len(v))
		//	}
		//	fmt.Println("-------------------")
		//	fmt.Println("-------------------")
		//}

	})
}

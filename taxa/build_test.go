package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bitbucket.org/heindl/taxa/utils"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of sources", t, func() {

		//srcs, err := fetchTaxaSources(context.Background(), inaturalist.ParseStringIDs("56830")...)
		//if err != nil {
		//	panic(err)
		//}
		//
		//marshalledSources, err := json.Marshal(srcs)
		//if err != nil {
		//	panic(err)
		//}

		b, err := ioutil.ReadFile("/Users/m/Desktop/morchella.json")
		So(err, ShouldBeNil)

		sources := CanonicalNameSources{}
		So(json.Unmarshal(b, &sources), ShouldBeNil)

		fmt.Println(utils.JsonOrSpew(sources.GenerateNameResults()))

	})
}

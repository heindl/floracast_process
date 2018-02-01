package inaturalist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"bitbucket.org/heindl/taxa/store"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch inaturalist", t, func() {

		usages, err := FetchNameUsages(context.Background(), 56830)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(usages))

		fmt.Println(utils.JsonOrSpew(usages.Names()))
		fmt.Println(utils.JsonOrSpew(usages.TargetIDs(store.DataSourceTypeGBIF)))

	})

	Convey("should fetch occurrences", t, func() {

		occurrences, err := FetchOccurrences(context.Background(), 56830)
		So(err, ShouldBeNil)
		fmt.Println("Total", len(occurrences))

	})

}

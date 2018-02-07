package inaturalist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"time"
	"bitbucket.org/heindl/taxa/datasources"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch inaturalist", t, func() {

		usages, err := FetchNameUsages(context.Background(), 56830)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(usages))


		//fmt.Println(utils.JsonOrSpew(usages.Names()))
		//fmt.Println(utils.JsonOrSpew(usages.TargetIDs(store.DataSourceTypeGBIF)))

	})

	SkipConvey("should fetch occurrences", t, func() {

		occurrences, err := FetchOccurrences(context.Background(), datasources.DataSourceTargetID("58682"), utils.TimePtr(time.Now().Add(time.Hour * 24 * 60 * -1)))
		So(err, ShouldBeNil)
		So(occurrences.Count(), ShouldEqual, 24)

		occurrences, err = FetchOccurrences(context.Background(), datasources.DataSourceTargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(occurrences.Count(), ShouldEqual, 100)

	})

}

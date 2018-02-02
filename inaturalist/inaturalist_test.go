package inaturalist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"bitbucket.org/heindl/taxa/store"
	"time"
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

	Convey("should fetch occurrences", t, func() {

		occurrences, err := FetchOccurrences(context.Background(), store.DataSourceTargetID("58682"), utils.TimePtr(time.Now().Add(time.Hour * 24 * 60 * -1)))
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 24)

		occurrences, err = FetchOccurrences(context.Background(), store.DataSourceTargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 100)

	})

}

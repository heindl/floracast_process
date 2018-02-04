package mushroom_observer

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"bitbucket.org/heindl/taxa/store"
	"time"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should generate a list of NameUsageSources", t, func() {

		res, err := MatchCanonicalNames("cantharellus cibarius")
		So(err, ShouldBeNil)

		fmt.Println(utils.JsonOrSpew(res))

	})

	Convey("should fetch occurrences ", t, func() {

		res, err := FetchOccurrences(context.Background(), store.DataSourceTargetID("404"), nil)
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 5)

		res, err = FetchOccurrences(context.Background(), store.DataSourceTargetID("404"), utils.TimePtr(time.Date(2009, time.January, 1, 0,  0, 0, 0, time.UTC)))
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 4)

	})
}

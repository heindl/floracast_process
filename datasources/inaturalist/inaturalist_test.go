package inaturalist

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/process/utils"
	"time"
	"bitbucket.org/heindl/process/datasources"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch inaturalist", t, func() {

		usages, err := FetchNameUsages(context.Background(), []string{"56830"}, nil)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(usages))


		//fmt.Println(utils.JsonOrSpew(usages.Names()))
		//fmt.Println(utils.JsonOrSpew(usages.TargetIDs(store.TypeGBIF)))

	})

	Convey("should fetch photos", t, func() {
		p, err := FetchPhotos(context.Background(), datasources.TargetID("58682"))
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(p))
	})

	SkipConvey("should fetch occurrences", t, func() {

		occurrences, err := FetchOccurrences(context.Background(), datasources.TargetID("58682"), utils.TimePtr(time.Now().Add(time.Hour * 24 * 60 * -1)))
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 24)

		occurrences, err = FetchOccurrences(context.Background(), datasources.TargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 100)

	})

}

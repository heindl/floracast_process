package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestInaturalistFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch taxon_schemes", t, func() {
		schemes, err := taxonID(56830).fetchTaxonSchemes()
		So(err, ShouldBeNil)
		So(len(schemes), ShouldEqual, 1)
		So(schemes[0].TargetID, ShouldEqual, datasources.TargetID("2594601"))
		So(schemes[0].SourceType, ShouldEqual, datasources.SourceType("27"))

		schemes, err = taxonID(96710).fetchTaxonSchemes()
		So(err, ShouldBeNil)
		So(len(schemes), ShouldEqual, 3)
		So(schemes[2].TargetID, ShouldEqual, datasources.TargetID("5184831"))
		So(schemes[2].SourceType, ShouldEqual, datasources.SourceType("27"))
		So(schemes[1].TargetID, ShouldEqual, datasources.TargetID("133384"))
		So(schemes[1].SourceType, ShouldEqual, datasources.SourceType("26"))
		So(schemes[0].TargetID, ShouldEqual, datasources.TargetID("133384"))
		So(schemes[0].SourceType, ShouldEqual, datasources.SourceType("14"))
	})

	Convey("should fetch inaturalist", t, func() {

		usages, err := FetchNameUsages(context.Background(), nil, datasources.TargetIDs{datasources.TargetID("56830")})
		So(err, ShouldBeNil)
		So(len(usages), ShouldEqual, 28)

	})

	Convey("should fetch photos", t, func() {
		photos, err := FetchPhotos(context.Background(), datasources.TargetID("58682"))
		So(err, ShouldBeNil)
		So(len(photos), ShouldEqual, 1)
		So(photos[0].Citation(), ShouldEqual, "(c) Leo Papandreou, some rights reserved (CC BY-NC-SA)")
		So(photos[0].Large(), ShouldEqual, "https://farm5.staticflickr.com/4031/4710142661_38bb26fb1a_b.jpg")
		So(photos[0].Thumbnail(), ShouldEqual, "https://farm5.staticflickr.com/4031/4710142661_38bb26fb1a_m.jpg")
		So(photos[0].SourceType(), ShouldEqual, datasources.TypeINaturalist)
	})

	Convey("should fetch occurrences", t, func() {
		// TODO: This test has a since time, so occurrence count may go up.

		baseTime, err := time.Parse("2006-01-02", "2017-12-29")
		So(err, ShouldBeNil)

		occurrences, err := FetchOccurrences(context.Background(), datasources.TargetID("58682"), utils.TimePtr(baseTime))
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 11)

		occurrences, err = FetchOccurrences(context.Background(), datasources.TargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(len(occurrences), ShouldEqual, 129)

	})

}

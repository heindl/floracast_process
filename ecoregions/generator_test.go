package ecoregions

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEcoRegionFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {

		cache, err := NewEcoRegionsCache()
		So(err, ShouldBeNil)

		r, err := cache.HasPoint(38.6270025, -90.1994042)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		So(r.EcoName, ShouldEqual, "Central forest-grasslands transition")
		So(r.EcoID, ShouldEqual, "245")

		r, err = cache.HasPoint(33.7676338, -84.5606888)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		So(r.EcoID, ShouldEqual, "70")
		So(r.EcoName, ShouldEqual, "Southeastern mixed forests")

		// St. Louis
		r, err = cache.HasPoint(38.6530169, -90.3835463)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
		So(r.EcoID, ShouldEqual, "245")
		So(r.EcoName, ShouldEqual, "Central forest-grasslands transition")

	})
}

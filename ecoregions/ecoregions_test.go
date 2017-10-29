package ecoregions

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestElevationFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {

		r, err := NewEcoRegionCache("/Users/m/Downloads/wwf_terr_ecos_oRn.json")
		So(err, ShouldBeNil)

		name, key, err := r.PointWithin(38.6270025, -90.1994042)
		So(err, ShouldBeNil)
		So(name, ShouldEqual, "Central forest-grasslands transition")
		So(key, ShouldEqual, "245")

		name, key, err = r.PointWithin(33.7676338,-84.5606888)
		So(err, ShouldBeNil)
		So(key, ShouldEqual, "70")
		So(name, ShouldEqual, "Southeastern mixed forests")

		// St. Louis
		name, key, err = r.PointWithin(38.6530169,-90.3835463)
		So(err, ShouldBeNil)
		So(key, ShouldEqual, "245")
		So(name, ShouldEqual, "Central forest-grasslands transition")


	})
}

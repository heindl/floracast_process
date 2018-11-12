package cache

import (
	"github.com/heindl/floracast_process/terra/ecoregions"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEcoRegionFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch ecoregions from cache", t, func() {

		region, err := FetchEcologicalRegion(38.6270025, -90.1994042)
		So(err, ShouldBeNil)
		So(region.Name(), ShouldEqual, "Central forest-grasslands transition")
		So(region.Realm(), ShouldEqual, 5)
		So(region.Biome(), ShouldEqual, 8)
		So(region.EcoNum(), ShouldEqual, 4)

		region, err = FetchEcologicalRegion(33.7676338, -84.5606888)
		So(err, ShouldBeNil)
		So(region.Name(), ShouldEqual, "Southeastern mixed forests")
		So(region.Realm(), ShouldEqual, 5)
		So(region.Biome(), ShouldEqual, 4)
		So(region.EcoNum(), ShouldEqual, 13)

		// St. Louis
		region, err = FetchEcologicalRegion(38.6530169, -90.3835463)
		So(err, ShouldBeNil)
		So(region.Name(), ShouldEqual, "Central forest-grasslands transition")
		So(region.Realm(), ShouldEqual, 5)
		So(region.Biome(), ShouldEqual, 8)
		So(region.EcoNum(), ShouldEqual, 4)

		_, err = FetchEcologicalRegion(37.867821, -122.3096932)
		So(err, ShouldEqual, ecoregions.ErrNotFound)

		region, err = FetchEcologicalRegion(19.3103, -99.3267)
		So(err, ShouldBeNil)
		So(region.Name(), ShouldEqual, "Trans-Mexican Volcanic Belt pine-oak forests")
		So(region.Realm(), ShouldEqual, 6)
		So(region.Biome(), ShouldEqual, 3)
		So(region.EcoNum(), ShouldEqual, 10)

	})
}

package ecoregions

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEcoRegionFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {

		cache, err := NewEcoRegionsCache()
		So(err, ShouldBeNil)

		//id := cache.EcoID(38.6270025, -90.1994042)
		//So(id.Valid(), ShouldBeTrue)
		//So(id.Name(), ShouldEqual, "Central forest-grasslands transition")
		//fmt.Println(id, id.Biome(), id.EcoNum())
		//
		//id = cache.EcoID(33.7676338, -84.5606888)
		//So(id.Valid(), ShouldBeTrue)
		//So(id.Name(), ShouldEqual, "Southeastern mixed forests")
		//fmt.Println(id, id.Biome(), id.EcoNum())
		//
		//// St. Louis
		//id = cache.EcoID(38.6530169, -90.3835463)
		//So(id.Valid(), ShouldBeTrue)
		//So(id.Name(), ShouldEqual, "Central forest-grasslands transition")
		//fmt.Println(id, id.Biome(), id.EcoNum())


		id := cache.EcoID(37.867821,-122.3096932)
		So(id.Valid(), ShouldBeTrue)
		fmt.Println(id.Name(), id.Biome(), id.EcoNum())

	})
}

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	speciesstore "bitbucket.org/heindl/species/store"
	"cloud.google.com/go/datastore"
)

func TestElevationFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {
		list := speciesstore.Occurrences{
			{Location: &datastore.GeoPoint{25.75027, -80.766463}},
			{Location: &datastore.GeoPoint{36.2340512,-116.8863299}},
			{Location: &datastore.GeoPoint{33.7676338,-84.5606914}},
		}
		So(setElevations(list), ShouldBeNil)

		So(list[0].Elevation, ShouldEqual, 4.377089023590088)
		So(list[1].Elevation, ShouldEqual, -59.19219207763672)
		So(list[2].Elevation, ShouldEqual, 241.1528778076172)
	})
}

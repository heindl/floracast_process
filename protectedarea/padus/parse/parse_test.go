package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProtectedAreaDatabaseParser(t *testing.T) {

	SkipConvey("should parse protected area json file", t, func() {

		processor, err := NewProcessor("/tmp/gap_analysis/GA/state.geojson", "/tmp/gap_analysis/GA/areas/")
		So(err, ShouldBeNil)

		collections, metrics, err := processor.ProcessFeatureCollections()
		So(err, ShouldBeNil)

		So(metrics.Total, ShouldEqual, 4230)
		So(metrics.Stats["Initial Filtered Total"], ShouldEqual, 1664)

		So(metrics.PublicAccessClosed, ShouldEqual, 910)
		So(metrics.PublicAccessUnknown, ShouldEqual, 671)
		So(metrics.PublicAccessRestricted, ShouldEqual, 49)
		So(metrics.EmptyAreas, ShouldEqual, 0)

		So(metrics.Stats["After Name Filter"], ShouldEqual, 1358)
		So(metrics.Stats["After Name Group"], ShouldEqual, 1173)
		So(metrics.Stats["After Centroid Distance Filter"], ShouldEqual, 1231)

		So(metrics.Stats["After Minimum Area Filter"], ShouldEqual, 502)
		So(metrics.Stats["After Cluster Decimation"], ShouldEqual, 216)

		So(len(collections), ShouldEqual, 216)

	})

	Convey("should parse protected area json file", t, func() {

		processor, err := NewProcessor("/tmp/gap_analysis/OR/state.geojson", "/tmp/gap_analysis/OR/areas/")
		So(err, ShouldBeNil)

		collections, metrics, err := processor.ProcessFeatureCollections()
		So(err, ShouldBeNil)

		So(metrics.Total, ShouldEqual, 7079)
		So(metrics.Stats["Initial Filtered Total"], ShouldEqual, 5360)

		So(metrics.PublicAccessClosed, ShouldEqual, 432)
		So(metrics.PublicAccessUnknown, ShouldEqual, 217)
		So(metrics.PublicAccessRestricted, ShouldEqual, 2141)
		So(metrics.EmptyAreas, ShouldEqual, 1)

		So(metrics.Stats["After Name Filter"], ShouldEqual, 4429)
		So(metrics.Stats["After Name Group"], ShouldEqual, 3690)
		So(metrics.Stats["After Centroid Distance Filter"], ShouldEqual, 4247)

		So(metrics.Stats["After Minimum Area Filter"], ShouldEqual, 1085)
		So(metrics.Stats["After Cluster Decimation"], ShouldEqual, 377)

		So(len(collections), ShouldEqual, 377)

	})
}

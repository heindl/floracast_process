package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"strings"
)

func TestProtectedAreaDatabaseParser(t *testing.T) {

	Convey("Given an Alabama state.geojson file and output directory", t, func() {

		// Note: Need to launch parser in advance: ./parse.sh AL
		stateGeoJSONFile := "/tmp/gap_analysis/AL/state.geojson"
		tmpDir, err := ioutil.TempDir("/tmp", "gap_analysis_test_")
		So(err, ShouldBeNil)

		Convey("Should create a new PAD-US processor, parse and write FeatureCollections to output directory", func() {

			processor, err := NewProcessor(stateGeoJSONFile, tmpDir)
			So(err, ShouldBeNil)

			collections, metrics, err := processor.ProcessFeatureCollections()
			So(err, ShouldBeNil)

			So(metrics.Total, ShouldEqual, 1083)
			So(metrics.Stats["Initial Filtered Total"], ShouldEqual, 520)

			So(metrics.PublicAccessClosed, ShouldEqual, 209)
			So(metrics.PublicAccessUnknown, ShouldEqual, 119)
			So(metrics.PublicAccessRestricted, ShouldEqual, 163)
			So(metrics.EmptyAreas, ShouldEqual, 0)

			So(metrics.Stats["After Name Filter"], ShouldEqual, 459)
			So(metrics.Stats["After Name Group"], ShouldEqual, 410)
			So(metrics.Stats["After Centroid Distance Filter"], ShouldEqual, 437)

			So(metrics.Stats["After Minimum Area Filter"], ShouldEqual, 259)
			So(metrics.Stats["After Cluster Decimation"], ShouldEqual, 110)

			So(len(collections), ShouldEqual, 110)

			So(processor.WriteCollections(collections), ShouldBeNil)

			// Ensure all files are accounted for and not empty.
			files, err := ioutil.ReadDir(tmpDir)
			So(err, ShouldBeNil)
			So(len(files), ShouldEqual, 110)
			for _, f := range files {
				So(f.Size(), ShouldBeGreaterThan, 0)
			}
		})

		Reset(func() {
			// This reset is run after each `Convey` at the same scope.
			So(strings.HasPrefix(tmpDir, "/tmp"), ShouldBeTrue)
			So(os.RemoveAll(tmpDir), ShouldBeNil)
		})

	})
}

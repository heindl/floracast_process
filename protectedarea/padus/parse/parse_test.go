package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProtectedAreaDatabaseParser(t *testing.T) {

	t.Parallel()

	Convey("should parse protected area json file", t, func() {

		processor, err := NewProcessor("/tmp", "/tmp")
		So(err, ShouldBeNil)

		collections, metrics, err := processor.ProcessFeatureCollections()
		So(err, ShouldBeNil)
		So(len(collections), ShouldEqual, 20)
		So(metrics.Total, ShouldEqual, 40)

	})
}

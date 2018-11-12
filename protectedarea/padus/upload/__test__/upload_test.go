package upload_test

import (
	"testing"

	"github.com/heindl/floracast_process/protectedarea"
	"github.com/heindl/floracast_process/store"
	"context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProtectedAreaDirectoryParser(t *testing.T) {

	Convey("Given a FloraStore and directory of ProtectedArea GeoJSON files", t, func() {

		testFiles := "./testFiles/"

		ctx := context.Background()

		floraStore, err := store.NewTestFloraStore(ctx)
		So(err, ShouldBeNil)

		protectedAreasCollection, err := floraStore.FirestoreCollection(store.CollectionProtectedAreas)
		So(err, ShouldBeNil)

		Convey("Should parse ProtectedArea directory for FeatureCollections", func() {

			areas, err := ParseProtectedAreaDirectory(testFiles)
			So(err, ShouldBeNil)
			So(len(areas), ShouldEqual, 5)

			for _, a := range areas {
				So(a.Valid(), ShouldBeTrue)
			}

			totalUploaded, err := areas.Upload(ctx, floraStore)
			So(err, ShouldBeNil)
			So(totalUploaded, ShouldEqual, 5)

			areaCount, err := floraStore.CountTestCollection(ctx, protectedAreasCollection)
			So(err, ShouldBeNil)
			So(areaCount, ShouldEqual, 5)

			// All should exist
			for _, queryArea := range areas {
				queryID, err := queryArea.ID()
				So(err, ShouldBeNil)
				fetchedArea, err := protectedarea.FetchOne(ctx, floraStore, queryID)
				So(err, ShouldBeNil)
				So(fetchedArea.Valid(), ShouldBeTrue)
			}

		})

		Reset(func() {
			So(floraStore.ClearTestCollection(ctx, protectedAreasCollection), ShouldBeNil)
			// Just to ensure our reset is successful ...
			areaCount, err := floraStore.CountTestCollection(ctx, protectedAreasCollection)
			So(err, ShouldBeNil)
			So(areaCount, ShouldEqual, 0)
		})

	})

}

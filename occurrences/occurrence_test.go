package occurrences

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/store"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	// These numbers may change because there is no end date.

	Convey("Should Fetch Occurrences for Taxon", t, func() {

		iNatAggr, err := FetchOccurrences(context.Background(), datasources.TypeINaturalist, datasources.TargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(iNatAggr.Count(), ShouldEqual, 121)

		gbifAggr, err := FetchOccurrences(context.Background(), datasources.TypeGBIF, datasources.TargetID("2594602"), nil)
		So(err, ShouldBeNil)
		So(gbifAggr.Count(), ShouldEqual, 205)

		So(iNatAggr.Merge(gbifAggr), ShouldBeNil)

		So(iNatAggr.Count(), ShouldEqual, 239)

	})

	Convey("Should Fetch Occurrences and Upload to FireStore", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewTestFloraStore(ctx)
		So(err, ShouldBeNil)

		occurrenceCollectionRef, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
		So(err, ShouldBeNil)

		Convey("Should generate a list of Random points and upload to FireStore", func() {

			aggr, err := FetchOccurrences(context.Background(), datasources.TypeINaturalist, datasources.TargetID("58682"), nil)
			So(err, ShouldBeNil)
			So(aggr.Count(), ShouldEqual, 121)

			So(aggr.Upload(context.Background(), floraStore), ShouldBeNil)

			floraStoreOccurrenceCount, err := floraStore.CountTestCollection(ctx, occurrenceCollectionRef)
			So(err, ShouldBeNil)
			So(floraStoreOccurrenceCount, ShouldEqual, 121)

		})

		Reset(func() {
			So(floraStore.ClearTestCollection(ctx, occurrenceCollectionRef), ShouldBeNil)
		})

	})
}

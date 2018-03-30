package occurrence

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	// These numbers may change because there is no end date.

	Convey("Should Fetch Occurrences for Taxon", t, func() {

		nameusageId, err := nameusage.NewNameUsageID()
		So(err, ShouldBeNil)

		iNatAggr, err := fetchOccurrencesForTarget(context.Background(), nameusageId, datasources.TypeINaturalist, datasources.TargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(iNatAggr.Count(), ShouldEqual, 120)

		gbifAggr, err := fetchOccurrencesForTarget(context.Background(), nameusageId, datasources.TypeGBIF, datasources.TargetID("2594602"), nil)
		So(err, ShouldBeNil)
		So(gbifAggr.Count(), ShouldEqual, 191)

		So(iNatAggr.Merge(gbifAggr), ShouldBeNil)

		So(iNatAggr.Count(), ShouldEqual, 223)

	})

	Convey("Should Fetch Occurrences and Upload to FireStore", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewTestFloraStore(ctx)
		So(err, ShouldBeNil)

		occurrenceCollectionRef, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
		So(err, ShouldBeNil)

		nameusageId, err := nameusage.NewNameUsageID()
		So(err, ShouldBeNil)

		Convey("Should generate a list of Random points and upload to FireStore", func() {

			aggr, err := fetchOccurrencesForTarget(context.Background(), nameusageId, datasources.TypeINaturalist, datasources.TargetID("58682"), nil)
			So(err, ShouldBeNil)
			So(aggr.Count(), ShouldEqual, 120)

			So(aggr.Upload(context.Background(), floraStore), ShouldBeNil)

			floraStoreOccurrenceCount, err := floraStore.CountTestCollection(ctx, occurrenceCollectionRef)
			So(err, ShouldBeNil)
			So(floraStoreOccurrenceCount, ShouldEqual, 120)

		})

		Reset(func() {
			So(floraStore.ClearTestCollection(ctx, occurrenceCollectionRef), ShouldBeNil)
		})

	})
}

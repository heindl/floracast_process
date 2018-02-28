package occurrences

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/store"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences for taxon", t, func() {

		aggr, err := FetchOccurrences(context.Background(), datasources.TypeINaturalist, datasources.TargetID("58682"), nil)
		So(err, ShouldBeNil)
		So(aggr.Count(), ShouldEqual, 120)

		gbifAggr, err := FetchOccurrences(context.Background(), datasources.TypeGBIF, datasources.TargetID("2594602"), nil)
		So(err, ShouldBeNil)
		So(gbifAggr.Count(), ShouldEqual, 205)

		So(aggr.Merge(gbifAggr), ShouldBeNil)

		So(aggr.Count(), ShouldEqual, 238)

		cxt := context.Background()

		florastore, err := store.NewTestFloraStore(cxt)
		So(err, ShouldBeNil)

		So(aggr.Upload(cxt, florastore), ShouldBeNil)

		col, err := florastore.FirestoreCollection(store.CollectionOccurrences)
		So(err, ShouldBeNil)

		snaps, err := col.Documents(cxt).GetAll()
		So(err, ShouldBeNil)
		So(len(snaps), ShouldEqual, 238)

		// TODO: Test duplication avoidance here.

		//for _, s := range snaps {
		//	_, err := s.Ref.Delete(cxt)
		//	So(err, ShouldBeNil)
		//}

	})
}

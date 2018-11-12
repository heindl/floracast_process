package generate

import (
	"github.com/heindl/floracast_process/store"
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	flag.Set("logtostderr", "true")
	flag.Set("v", "0")
	flag.Parse()

	// These numbers may change because there is no end date.

	//SkipConvey("Should download saved model", t, func() {
	//
	//	ctx := context.Background()
	//
	//	floraStore, err := store.NewFloraStore(ctx)
	//	So(err, ShouldBeNil)
	//
	//	mdllr, err := NewModeller(floraStore)
	//	So(err, ShouldBeNil)
	//
	//	_, err = mdllr.FetchModel(ctx, "9sykdre6ougztwabsjjufiwvu")
	//	So(err, ShouldBeNil)
	//
	//	So(mdllr.Close(), ShouldBeNil)
	//
	//})

	Convey("Should fetch occurrences", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewFloraStore(ctx)
		So(err, ShouldBeNil)

		collection, err := GeneratePredictions(
			ctx,
			"qWlT2bh",
			floraStore,
			nil,
			"/tmp/uJjIDtxo/model/exports/1524784583",
			"/Users/m/Desktop/protected_areas/**/*.tfrecords",
		)
		So(err, ShouldBeNil)
		So(collection.Count(), ShouldEqual, 68936)

	})
}

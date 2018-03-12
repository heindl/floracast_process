package generate

import (
	"bitbucket.org/heindl/process/store"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	// These numbers may change because there is no end date.

	SkipConvey("Should download saved model", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewFloraStore(ctx)
		So(err, ShouldBeNil)

		mdllr, err := NewModeller(floraStore)
		So(err, ShouldBeNil)

		_, err = mdllr.FetchModel(ctx, "9sykdre6ougztwabsjjufiwvu")
		So(err, ShouldBeNil)

		So(mdllr.Close(), ShouldBeNil)

	})

	Convey("Should fetch occurrences", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewFloraStore(ctx)
		So(err, ShouldBeNil)

		list, err := GeneratePredictions(ctx, "9sykdre6ougztwabsjjufiwvu", floraStore, &DateRange{
			Start: "20180201",
			End:   "20180312",
		})
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 707)

	})
}

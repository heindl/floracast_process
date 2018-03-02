package occurrence

import (
	"bitbucket.org/heindl/process/store"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"
)

func TestRandomProvider(t *testing.T) {

	Convey("Given a FireStore Connection", t, func() {

		ctx := context.Background()

		floraStore, err := store.NewTestFloraStore(ctx)
		So(err, ShouldBeNil)

		randomCollectionRef, err := floraStore.FirestoreCollection(store.CollectionRandom)
		So(err, ShouldBeNil)

		Convey("Should generate a list of Random points and upload to FireStore", func() {

			gridLevelFour, err := GenerateRandomOccurrences(4, 1)
			So(err, ShouldBeNil)
			So(gridLevelFour.Count(), ShouldEqual, 216)

			gridLevelTwo, err := GenerateRandomOccurrences(3, 1)
			So(err, ShouldBeNil)
			So(gridLevelTwo.Count(), ShouldEqual, 64)

			So(gridLevelTwo.Upload(context.Background(), floraStore), ShouldBeNil)

			floraStoreRandomCount, err := floraStore.CountTestCollection(ctx, randomCollectionRef)
			So(err, ShouldBeNil)
			So(floraStoreRandomCount, ShouldEqual, 64)

		})

		Reset(func() {
			So(ClearRandomPoints(ctx, floraStore), ShouldBeNil)
			floraStoreRandomCount, err := floraStore.CountTestCollection(ctx, randomCollectionRef)
			So(err, ShouldBeNil)
			So(floraStoreRandomCount, ShouldEqual, 0)
		})

	})
}

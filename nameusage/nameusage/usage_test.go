package nameusage

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/store"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"testing"
)

func TestNameUsage(t *testing.T) {

	t.Parallel()

	Convey("Should Create, Parse and Upload NameUsage", t, func() {

		ctx := context.Background()
		floraStore, err := store.NewTestFloraStore(ctx)
		So(err, ShouldBeNil)
		nameUsageCollection, err := floraStore.FirestoreCollection(store.CollectionNameUsages)
		So(err, ShouldBeNil)

		Convey("Should Create, Parse and Upload NameUsage", func() {

			targetID, err := datasources.NewDataSourceTargetIDFromInt(datasources.TypeGBIF, 12345)
			So(err, ShouldBeNil)

			name, err := canonicalname.NewCanonicalName("morchella esculenta", "species")
			So(err, ShouldBeNil)

			src, err := NewSource(datasources.TypeGBIF, targetID, name)
			So(err, ShouldBeNil)

			So(src.RegisterOccurrenceFetch(5), ShouldBeNil)

			initialUsage, err := NewNameUsage(src)
			So(err, ShouldBeNil)

			initialID, err := initialUsage.ID()
			So(err, ShouldBeNil)

			deletedIDs, err := initialUsage.Upload(ctx, floraStore)
			So(err, ShouldBeNil)
			So(len(deletedIDs), ShouldEqual, 0)

			docCount, err := floraStore.CountTestCollection(ctx, nameUsageCollection)
			So(err, ShouldBeNil)
			So(docCount, ShouldEqual, 1)

			iter := FetchAll(ctx, floraStore)

			fetchedUsage, err := iter.Next()
			So(err, ShouldBeNil)

			_, err = iter.Next()
			So(err, ShouldEqual, iterator.Done)

			So(initialUsage.CanonicalName().Equals(fetchedUsage.CanonicalName()), ShouldBeTrue)

			fetchedID, err := fetchedUsage.ID()
			So(err, ShouldBeNil)

			So(initialID, ShouldEqual, fetchedID)

		})

		Reset(func() {
			floraStore.ClearTestCollection(ctx, nameUsageCollection)
		})

	})
}

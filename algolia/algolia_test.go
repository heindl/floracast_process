package algolia

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/processors/utils"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/processors/store"
	"time"
)

func TestNameUsageProcessor(t *testing.T) {

	t.Parallel()

	Convey("should generate occurrence usages", t, func() {

		id1 := nameusage.NameUsageID("eM3R8X2YQyLJWiLMVIGzZaU1I")
		id2 := nameusage.NameUsageID("aM3R8X2YQyLJWiLMVIGzZaU1I")

		usage, err := nameusage.NameUsageFromJSON(id1, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)

		objs, err := generateNameUsageObjects(context.Background(), usage)
		So(err, ShouldBeNil)
		So(len(objs), ShouldEqual, 103)

		batches := objs.batches(20)
		So(len(batches), ShouldEqual, 6)

		florastore, err := store.NewTestFloraStore(context.Background())
		So(err, ShouldBeNil)

		So(uploadNameUsageObjects(florastore, objs), ShouldBeNil)

		time.Sleep(time.Second * 5) // Consistency

		c, err := countNameUsages(florastore, id1)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, 102)

		usage, err = nameusage.NameUsageFromJSON(id2, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)

		So(UploadNameUsageObjects(context.Background(), florastore, usage, id1), ShouldBeNil)

		time.Sleep(time.Second * 5) // Consistency

		c, err = countNameUsages(florastore, id1, id2)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, 102)

		time.Sleep(time.Second * 5) // Consistency

		So(deleteNameUsageObjects(florastore, id1, id2), ShouldBeNil)

		time.Sleep(time.Second * 5) // Consistency

		c, err = countNameUsages(florastore, id1, id2)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, 0)

	})
}

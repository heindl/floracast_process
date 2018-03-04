package algolia

import (
	"testing"
	"time"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/dropbox/godropbox/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNameUsageProcessor(t *testing.T) {

	t.Parallel()

	Convey("should generate occurrence usages", t, func() {

		id1 := nameusage.ID("eM3R8X2YQyLJWiLMVIGzZaU1I")
		id2 := nameusage.ID("aM3R8X2YQyLJWiLMVIGzZaU1I")

		usage, err := nameusage.FromJSON(id1, utils.GetFetchedMorchellaUsageTestData())
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
		So(c, ShouldEqual, 103)

		usage, err = nameusage.FromJSON(id2, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)

		So(UploadNameUsageObjects(context.Background(), florastore, usage, id1), ShouldBeNil)

		time.Sleep(time.Second * 5) // Consistency

		c, err = countNameUsages(florastore, id1, id2)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, 103)

		time.Sleep(time.Second * 5) // Consistency

		So(deleteNameUsageObjects(florastore, id1, id2), ShouldBeNil)

		time.Sleep(time.Second * 5) // Consistency

		c, err = countNameUsages(florastore, id1, id2)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, 0)

	})
}

func countNameUsages(florastore store.FloraStore, nameUsageIDs ...nameusage.ID) (int, error) {
	index, err := florastore.AlgoliaIndex(nameUsageIndex)
	if err != nil {
		return 0, err
	}
	count := 0

	iter, err := index.BrowseAll(nameUsageIDFilter(nameUsageIDs...))
	if err != nil && err != algoliasearch.NoMoreHitsErr {
		return 0, errors.Wrap(err, "Could not browse Algolia NameUsage Index")
	}

	for {
		if _, err := iter.Next(); err != nil && err != algoliasearch.NoMoreHitsErr {
			return 0, err
		} else if err != nil && err == algoliasearch.NoMoreHitsErr {
			break
		}
		count++
	}
	return count, nil

}

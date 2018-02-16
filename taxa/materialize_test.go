package taxa

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"bitbucket.org/heindl/processors/utils"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/processors/store"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch description", t, func() {
		id1 := nameusage.NameUsageID("eM3R8X2YQyLJWiLMVIGzZaU1I")
		usage, err := nameusage.NameUsageFromJSON(id1, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)
		desc, err := fetchDescription(context.Background(), usage)
		So(err, ShouldBeNil)
		So(desc, ShouldNotBeNil)
		So(desc.Citation, ShouldEqual, `Wikipedia contributors. "Morchella esculenta". Wikipedia, The Free Encyclopedia. 5 Feb. 2018. Web. 14 Feb. 2018. <http://en.wikipedia.org/wiki/Morchella_esculenta>`)
		So(desc.Text, ShouldEqual, "Morchella esculenta, (commonly known as common morel, morel, yellow morel, true morel, morel mushroom, and sponge morel) is a species of fungus in the Morchellaceae family of the Ascomycota. It is one of the most readily recognized of all the edible mushrooms and highly sought after. Each fruit body begins as a tightly compressed, grayish sponge with lighter ridges, and expands to form a large yellowish sponge with large pits and ridges raised on a large...")
	})

	Convey("should materialize name usage", t, func() {

		id1 := nameusage.NameUsageID("eM3R8X2YQyLJWiLMVIGzZaU1I")

		usage, err := nameusage.NameUsageFromJSON(id1, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)

		m, err := materialize(context.Background(), usage)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

	})

	Convey("should upload materialized name usage", t, func() {

		id1 := nameusage.NameUsageID("eM3R8X2YQyLJWiLMVIGzZaU1I")

		usage, err := nameusage.NameUsageFromJSON(id1, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)

		cxt := context.Background()

		florastore, err := store.NewTestFloraStore(cxt)
		So(err, ShouldBeNil)

		So(UploadMaterializedTaxa(cxt, florastore, usage), ShouldBeNil)

		id2 := nameusage.NameUsageID("aM3R8X2YQyLJWiLMVIGzZaU1I")

		usage, err = nameusage.NameUsageFromJSON(id2, utils.GetFetchedMorchellaUsageTestData())
		So(err, ShouldBeNil)
		So(UploadMaterializedTaxa(cxt, florastore, usage, id1), ShouldBeNil)

		So(clearMaterializedTaxa(cxt, florastore, nameusage.NameUsageIDs{id2}), ShouldBeNil)

	})
}
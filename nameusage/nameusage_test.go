package algolia

import (
	"github.com/heindl/floracast_process/algolia"
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/nameusage/canonicalname"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/taxa"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMaterialization(t *testing.T) {

	t.Parallel()

	Convey("Test Materialization", t, func() {

		targetID, err := datasources.NewTargetID("58682", datasources.TypeINaturalist)
		So(err, ShouldBeNil)

		src, err := nameusage.NewSource(
			datasources.TypeINaturalist,
			targetID,
			&canonicalname.Name{SciName: "Morchella esculenta"},
		)
		So(err, ShouldBeNil)

		So(src.AddCommonNames("Morel"), ShouldBeNil)

		So(src.AddSynonym(&canonicalname.Name{SciName: "Helvella esculenta"}), ShouldBeNil)

		So(src.AddSynonym(&canonicalname.Name{SciName: "Morchella abientina"}), ShouldBeNil)

		So(src.RegisterOccurrenceFetch(5), ShouldBeNil)

		usage, err := nameusage.NewNameUsage(src)
		So(err, ShouldBeNil)

		cxt := context.Background()

		floraStore, err := store.NewTestFloraStore(cxt)
		So(err, ShouldBeNil)

		deletedNameUsages, err := usage.Upload(cxt, floraStore)
		So(err, ShouldBeNil)

		So(taxa.UploadMaterializedTaxon(cxt, floraStore, usage, deletedNameUsages...), ShouldBeNil)

		id, err := usage.ID()
		So(err, ShouldBeNil)

		So(algolia.IndexNameUsage(cxt, floraStore, id), ShouldBeNil)

	})
}

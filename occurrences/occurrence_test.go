package occurrences

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/processors/utils"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/processors/store"
	"github.com/mongodb/mongo-tools/common/json"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should upload materialized name usage", t, func() {

		b := utils.GetFileContents("./testdata.json")

		oa := OccurrenceAggregation{}
		So(json.Unmarshal(b, &oa), ShouldBeNil)

		So(oa.Count(), ShouldEqual, 238)

		cxt := context.Background()

		florastore, err := store.NewTestFloraStore(cxt)
		So(err, ShouldBeNil)

		So(oa.Upload(cxt, florastore), ShouldBeNil)

		snaps, err := florastore.FirestoreCollection(store.CollectionOccurrences).Documents(cxt).GetAll()
		So(err, ShouldBeNil)
		So(len(snaps), ShouldEqual, 238)

		So(oa.Upload(cxt, florastore), ShouldBeNil)

		snaps, err = florastore.FirestoreCollection(store.CollectionOccurrences).Documents(cxt).GetAll()
		So(err, ShouldBeNil)
		So(len(snaps), ShouldEqual, 238)

		for _, s := range snaps {
			_, err := s.Ref.Delete(cxt)
			So(err, ShouldBeNil)
		}

	})
}
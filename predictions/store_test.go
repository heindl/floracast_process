package predictions

import (
	"testing"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"math/rand"
)

func TestPredictionUpload(t *testing.T) {

	Convey("given predictions, should upload into GeoHashIndex", t, func() {

		cxt := context.Background()

		floraStore, err := store.NewFloraStore(cxt)
		So(err, ShouldBeNil)

		for i := 0; i <= 3; i++ {
			id, err := nameusage.NewNameUsageID()
			So(err, ShouldBeNil)
			c, err := NewCollection(id, floraStore)
			So(err, ShouldBeNil)

			So(c.Add(35.06375602820245, -90.11736207549293, utils.FormattedDate("20180102"), rand.Float64()), ShouldBeNil)
			So(c.Add(35.06375602820245, -90.11736207549293, utils.FormattedDate("20180108"), rand.Float64()), ShouldBeNil)
			So(c.Add(35.06375602820245, -90.11736207549293, utils.FormattedDate("20180112"), rand.Float64()), ShouldBeNil)

			So(c.Add(40.17699813751862, -80.53083704621974, utils.FormattedDate("20180102"), rand.Float64()), ShouldBeNil)
			So(c.Add(40.17699813751862, -80.53083704621974, utils.FormattedDate("20180108"), rand.Float64()), ShouldBeNil)
			So(c.Add(40.17699813751862, -80.53083704621974, utils.FormattedDate("20180112"), rand.Float64()), ShouldBeNil)
			So(c.Upload(cxt), ShouldBeNil)
		}

	})
}

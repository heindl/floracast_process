package mushroom_observer

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of NameUsageSources", t, func() {

		res, err := MatchCanonicalNames("cantharellus cibarius")
		So(err, ShouldBeNil)

		fmt.Println(utils.JsonOrSpew(res))


	})
}

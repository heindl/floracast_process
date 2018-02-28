package natureserve

import (
	"bitbucket.org/heindl/process/datasources"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch natureserve", t, func() {

		names := []string{"morchella deliciosa", "cantharellus cibarius", "boletus edulis"}

		res, err := FetchNameUsages(context.Background(), names, nil)
		So(err, ShouldBeNil)

		for i, ex := range [][2]string{
			{"cantharellus cibarius", "ELEMENT_GLOBAL.2.122116"},
			{"boletus edulis", "ELEMENT_GLOBAL.2.878680"},
		} {
			names, err := res[i].AllScientificNames()
			So(err, ShouldBeNil)
			So(names, ShouldContain, ex[0])
			hasSrc, err := res[i].HasSource(datasources.TypeNatureServe, datasources.TargetID(ex[1]))
			So(err, ShouldBeNil)
			So(hasSrc, ShouldBeTrue)
		}
	})
}

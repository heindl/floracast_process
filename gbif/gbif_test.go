package gbif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch gbif", t, func() {
		taxa, err := MatchNames(context.Background(), "Morchella esculenta")
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(taxa))
	})
}

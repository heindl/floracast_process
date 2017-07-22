package occurrences

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	speciesstore "bitbucket.org/heindl/species/store"
	//"bitbucket.org/heindl/utils"
	"github.com/jonboulle/clockwork"
	"fmt"
)

func TestOccurrenceFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {

		taxastore := speciesstore.NewTestTaxaStore()

		fetcher := NewOccurrenceFetcher(taxastore, clockwork.NewFakeClockAt(time.Date(2017, time.May, 18, 0, 0, 0, 0, time.UTC)))

		So(fetcher.FetchOccurrences(), ShouldBeNil)

		ocs, err := taxastore.GetOccurrences(nil)
		So(err, ShouldBeNil)

		fmt.Println(len(ocs))

		//fmt.Println(utils.JsonOrSpew(ocs[0:20]))
	})
}

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	speciesstore "bitbucket.org/heindl/taxa/store"
	//"bitbucket.org/taxa/utils"
	"github.com/jonboulle/clockwork"
	"fmt"
	"bitbucket.org/taxa/utils"
)

func TestOccurrenceFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("show occurrences", t, func() {

		taxastore := speciesstore.NewTestTaxaStore()

		ocs, err := taxastore.GetOccurrences(nil)
		So(err, ShouldBeNil)

		for _, o := range ocs {
			if o.Elevation == 0 {
				fmt.Println(utils.JsonOrSpew(o))
			}
		}

		Reset(func() {
			So(taxastore.Close(), ShouldBeNil)
		})
	})

	Convey("should fetch occurrences and add to queue", t, func() {

		taxastore := speciesstore.NewTestTaxaStore()

		//taxa, err := taxastore.ReadTaxa()
		//So(err, ShouldBeNil)

		//for _, t := range taxa {
			list, err := taxastore.GetOccurrences(nil)
			So(err, ShouldBeNil)
			m := map[int64]int{}
			for _, l := range list {
				if _, ok := m[l.Key.Parent.Parent.ID]; !ok {
					m[l.Key.Parent.Parent.ID] = 1
				} else {
					m[l.Key.Parent.Parent.ID] += 1
				}
			}


			fmt.Println(m)

		Reset(func() {
			So(taxastore.Close(), ShouldBeNil)
		})
	})

	SkipConvey("should fetch occurrences and add to queue", t, func() {

		taxastore := speciesstore.NewTestTaxaStore()

		fetcher := NewOccurrenceFetcher(taxastore, clockwork.NewFakeClockAt(time.Date(2017, time.May, 18, 0, 0, 0, 0, time.UTC)))

		So(fetcher.FetchOccurrences(), ShouldBeNil)

		ocs, err := taxastore.GetOccurrences(nil)
		So(err, ShouldBeNil)

		fmt.Println(len(ocs))

		fmt.Println(utils.JsonOrSpew(ocs[0:20]))

		Reset(func() {
			So(taxastore.Close(), ShouldBeNil)
		})
	})
}

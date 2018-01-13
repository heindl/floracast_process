package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	//"bitbucket.org/taxa/utils"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"fmt"
	"github.com/jonboulle/clockwork"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func TestOccurrenceFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("show occurrences", t, func() {

		taxastore := store.NewTestTaxaStore()

		ocs, err := taxastore.GetOccurrences(context.Background(), store.TaxonID(""))
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

	SkipConvey("should fetch occurrences and add to queue", t, func() {

		taxastore := store.NewTestTaxaStore()

		//taxa, err := taxastore.ReadTaxa()
		//So(err, ShouldBeNil)

		//for _, t := range taxa {
		list, err := taxastore.GetOccurrences(context.Background(), store.TaxonID(""))
		So(err, ShouldBeNil)
		m := map[store.TaxonID]int{}
		for _, l := range list {
			if _, ok := m[l.TaxonID]; !ok {
				m[l.TaxonID] = 1
			} else {
				m[l.TaxonID] += 1
			}
		}

		fmt.Println(m)

		Reset(func() {
			So(taxastore.Close(), ShouldBeNil)
		})
	})

	Convey("should fetch occurrences and add to queue", t, func() {

		taxastore := store.NewTestTaxaStore()

		fetcher := NewOccurrenceFetcher(taxastore, clockwork.NewFakeClockAt(time.Date(2017, time.May, 18, 0, 0, 0, 0, time.UTC)))

		So(fetcher.FetchOccurrences(), ShouldBeNil)

		ocs, err := taxastore.GetOccurrences(context.Background(), store.TaxonID("58583"))
		So(err, ShouldBeNil)

		fmt.Println("occurrence_length", len(ocs))

		fmt.Println(utils.JsonOrSpew(ocs[0:20]))

		Reset(func() {
			So(taxastore.Close(), ShouldBeNil)
		})
	})
}

func TestElevationFetch(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences and add to queue", t, func() {
		list := store.Occurrences{
			{Location: latlng.LatLng{25.75027, -80.766463}},
			{Location: latlng.LatLng{36.2340512, -116.8863299}},
			{Location: latlng.LatLng{33.7676338, -84.5606914}},
		}
		So(setElevations(list), ShouldBeNil)

		So(list[0].Elevation, ShouldEqual, 4.377089023590088)
		So(list[1].Elevation, ShouldEqual, -59.19219207763672)
		So(list[2].Elevation, ShouldEqual, 241.1528778076172)
	})
}

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/taxa/store"
	"github.com/jonboulle/clockwork"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"golang.org/x/net/context"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch inaturalist schemes", t, func() {
		f := fetcher{}
		srcs, err := f.fetchDataSources(store.TaxonID("58583"), store.CanonicalName("Limenitis arthemis ssp. arthemis"), true)
		So(err, ShouldBeNil)
		So(len(srcs), ShouldEqual, 3)
		So(srcs[0].SourceID, ShouldEqual, store.DataSourceID("11"))
		So(srcs[0].Kind, ShouldEqual, store.DataSourceKindDescription)
		So(srcs[1].SourceID, ShouldEqual, store.DataSourceID("27"))
		So(srcs[1].Kind, ShouldEqual, store.DataSourceKindOccurrence)
		So(srcs[2].SourceID, ShouldEqual, store.DataSourceID("27"))
		So(srcs[2].Kind,ShouldEqual, store.DataSourceKindPhoto)
	})

	Convey("should fetch additional taxon ids from the gbif based on the canonical name", t, func() {
		f := fetcher{}
		ids, err := f.fetchAdditionalGBIFTaxonIDs("Morchella esculenta", store.DataSourceTargetID("2594602"))
		So(err, ShouldBeNil)
		So(ids[0], ShouldEqual, store.TaxonID("8574619"))

		ids, err = f.fetchAdditionalGBIFTaxonIDs("Cantharellus cibarius", store.DataSourceTargetID("5249504"))
		So(err, ShouldBeNil)
		Println(ids)
	})

	Convey("should fetch all species in subfamily Limenitidinae", t, func() {

		f := fetcher{
			Store: store.NewMockStore(t),
			Clock: clockwork.NewFakeClock(),
		}

		cxt := context.Background()

		taxa, err := f.Store.ReadTaxa(cxt)
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 0)

		So(f.FetchProcessTaxa(cxt, []store.TaxonID{store.TaxonID("58583")}), ShouldBeNil)

		taxa, err = f.Store.ReadTaxa(cxt)
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 19)

		taxa, err = f.Store.ReadSpecies(cxt)
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 7)

		dataSources, err := f.Store.GetOccurrenceDataSources(cxt, store.TaxonID(""))
		So(err, ShouldBeNil)
		So(len(dataSources), ShouldEqual, 6)
		have := []string{}
		for i := range dataSources {
			if utils.Contains(have, string(dataSources[i].TaxonID)) {
				fmt.Println(dataSources[i].TaxonID)
			}
			have = append(have, string(dataSources[i].TaxonID))
		}

		dataSources, err = f.Store.GetOccurrenceDataSources(cxt, taxa[0].ID)
		So(err, ShouldBeNil)
		So(len(dataSources), ShouldEqual, 1)

		dataSources, err = f.Store.GetOccurrenceDataSources(cxt, taxa[2].ID)
		So(err, ShouldBeNil)
		So(len(dataSources), ShouldEqual, 1)

		dataSources, err = f.Store.GetOccurrenceDataSources(cxt, taxa[3].ID)
		So(err, ShouldBeNil)
		So(len(dataSources), ShouldEqual, 1)

		Reset(func() {
			So(f.Store.Close(), ShouldBeNil)
		})

	})
}

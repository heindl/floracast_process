package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/species/store"
	"github.com/jonboulle/clockwork"
	"fmt"
	"bitbucket.org/heindl/utils"
	"strconv"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch inaturalist schemes", t, func() {
		f := fetcher{}
		schemes, err := f.fetchSchemes(store.NewTaxonKey(58585, store.RankSubSpecies), true)
		So(err, ShouldBeNil)
		So(len(schemes), ShouldEqual, 3)
		So(schemes[0].Key.Name, ShouldEqual, store.SchemeSourceID("11|||ELEMENT_GLOBAL.2.108251"))
		So(schemes[0].Key.Kind, ShouldEqual, store.EntityKindMetaScheme)
		So(schemes[1].Key.Name, ShouldEqual, store.SchemeSourceID("27|||5714327"))
		So(schemes[1].Key.Kind, ShouldEqual, store.EntityKindOccurrenceScheme)
		So(schemes[2].Key.Name, ShouldEqual, store.SchemeSourceID("27|||5714327"))
		So(schemes[2].Key.Kind,ShouldEqual, store.EntityKindMetaScheme)
	})

	Convey("should fetch all species in subfamily Limenitidinae", t, func() {

		f := fetcher{
			Store: store.NewMockStore(t),
			Clock: clockwork.NewFakeClock(),
		}

		taxa, err := f.Store.ReadTaxa()
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 0)

		So(f.FetchProcessTaxa(50881), ShouldBeNil)

		taxa, err = f.Store.ReadTaxa()
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 51)

		taxa, err = f.Store.ReadSpecies()
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 39)

		schema, err := f.Store.GetOccurrenceSchema(nil)
		So(err, ShouldBeNil)
		So(len(schema), ShouldEqual, 36)
		have := []string{}
		for i := range schema {
			if utils.Contains(have, strconv.Itoa(int(schema[i].Key.Parent.ID))) {
				fmt.Println(schema[i].Key.Parent.ID)
			}
			have = append(have, strconv.Itoa(int(schema[i].Key.Parent.ID)))
		}

		schema, err = f.Store.GetOccurrenceSchema(taxa[0].Key)
		So(err, ShouldBeNil)
		So(len(schema), ShouldEqual, 1)

		schema, err = f.Store.GetOccurrenceSchema(taxa[1].Key)
		So(err, ShouldBeNil)
		So(len(schema), ShouldEqual, 1)

		schema, err = f.Store.GetOccurrenceSchema(taxa[2].Key)
		So(err, ShouldBeNil)
		So(len(schema), ShouldEqual, 1)

		Reset(func() {
			So(f.Store.Close(), ShouldBeNil)
		})

	})
}

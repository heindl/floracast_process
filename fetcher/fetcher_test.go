package main

import (
	"bitbucket.org/heindl/nsqeco"
	"bitbucket.org/heindl/species/store"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/species"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("given an environment", t, func() {

		producer := &nsqeco.MockProducer{}
		store := store.NewMockStore(t)

		Convey("should handle taxon fetch message without error", func() {

			fetcher := &SpeciesFetcher{
				NSQProducer:  producer,
				SpeciesStore: store,
			}
			So(fetcher.FetchSpecies(species.CanonicalName("Limenitis arthemis")), ShouldBeNil)

			specs, err := store.Read()
			So(err, ShouldBeNil)

			So(len(specs), ShouldEqual, 8)
			for _, n := range []string{
				"Limenitis arthemis arizonensis",
				"Limenitis arthemis rubrofasciata",
				"Limenitis arthemis albofasciata",
				"Limenitis arthemis arthemis",
				"Limenitis arthemis astyanax",
			} {
				var exists bool
				for _, s := range specs {
					if n == string(s.CanonicalName) {
						exists = true
					}
				}
				So(exists, ShouldBeTrue)
			}
			// Check for EOL data
			So(len(specs), ShouldEqual, 8)
			So(specs[0].CanonicalName, ShouldEqual, species.CanonicalName("Limenitis arthemis"))
			So(len(specs[0].Sources), ShouldEqual, 9)
			So(specs[0].Image, ShouldNotBeNil)
			So(specs[0].Image.Value, ShouldEqual, "http://media.eol.org/content/2010/12/10/03/77851_orig.jpg")
			So(specs[0].Description, ShouldNotBeNil)
			So(specs[0].Description.Value, ShouldEqual, "Similar to our other two species of true admirals (\u003ci\u003eLimenitis\u003c/i\u003e). The White lacks the rust-coloured forewing tips of Lorquin's (\u003ci\u003eL. lorquini\u003c/i\u003e), and has a row of reddish spots bordering the outside of the hindwing white band. Their ranges overlap only in the Waterton - Crowsnest region, where hybrid individuals exhibiting characters intermediate between the White and Lorquin's are sometimes found. \u003ci\u003eL. arthemis\u003c/i\u003e also has more orange on the hindwing upperside than Weidemeyer's (\u003ci\u003eL. weidemeyerii\u003c/i\u003e), and has a red-brown hindwing underside base rather than predominantly white. Hybrids between these two species sometimes also occur. \u0026nbsp;The western Canadian populations are subspecies \u003ci\u003erubrofasciata\u003c/i\u003e.  \u0026nbsp;")

			//So(specs[0].CanonicalName, ShouldEqual, "Limenitis arthemis arizonensis")
			//So(len(specs[0].Sources), ShouldEqual, 1)
			//So(specs[0].Sources[0].IndexKey, ShouldEqual, 6225972)

		})


		Reset(func() {
			store.Close()
		})

	})


}

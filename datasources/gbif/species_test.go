package gbif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSpecies(t *testing.T) {

	t.Parallel()

	spcs := species(5231190)

	Convey("given an encyclopedia of life species", t, func() {

		Convey("should return a single name usage", func() {
			n, err := spcs.Name()
			So(err, ShouldBeNil)
			So(n.Key, ShouldEqual, 5231190)
			So(n.Kingdom, ShouldEqual, "Animalia")
			So(n.Phylum, ShouldEqual, "Chordata")
		})

		Convey("should return a parsed name usage", func() {
			n, err := spcs.ParsedName()
			So(err, ShouldBeNil)
			So(n.ScientificName, ShouldEqual, "Passer domesticus (Linnaeus, 1758)")
			So(n.Type, ShouldEqual, "SCIENTIFIC")
			So(n.AuthorsParsed, ShouldEqual, false)
		})

		Convey("should return a list of parent name usages", func() {
			r, err := spcs.Parents()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 6)
			So(r[0].ScientificName, ShouldEqual, "Animalia")
			So(r[0].KingdomKey, ShouldEqual, 1)
		})

		Convey("should return a list of child name usages", func() {
			r, err := spcs.Children()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 15)
		})

		Convey("should return a list of related name usages", func() {
			r, err := spcs.Related()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 43)
		})

		Convey("should return a list of synonymical name usages", func() {
			r, err := spcs.Synonyms()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 1)
		})

		Convey("should return a list of distributions", func() {
			r, err := spcs.fetchDistributions()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 16)
		})

		Convey("should return a list of media", func() {
			r, err := spcs.fetchMedia()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 0)
			//So(r[0].Type, ShouldEqual, "StillImage")
			//So(r[0].Identifier, ShouldEqual, "http://upload.wikimedia.org/wikipedia/commons/d/d9/House_sparrowII.jpg")
			//So(r[0].SourceTaxonKey, ShouldEqual, 100448484)
		})

		Convey("should return a list of references", func() {
			r, err := spcs.fetchReferences()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 13)
		})

		Convey("should return a list of vernacular names", func() {
			r, err := spcs.fetchVernacularNames()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 74)
			So(r[0].VernacularName, ShouldEqual, "English sparrow")
			So(r[0].Language, ShouldEqual, "")
			So(r[0].Source, ShouldEqual, "Global Invasive Species Database")
			So(r[0].SourceTaxonKey, ShouldEqual, 100220560)
		})

	})

}

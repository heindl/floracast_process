package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSpecies(t *testing.T) {

	t.Parallel()

	species := Species(5231190)

	Convey("given an encyclopedia of life species", t, func() {

		Convey("should return a single name usage", func() {
			n, err := species.Name()
			So(err, ShouldBeNil)
			So(n.TaxonID, ShouldEqual, "119127243")
			So(n.Kingdom, ShouldEqual, "Animalia")
			So(n.Phylum, ShouldEqual, "Chordata")
		})

		Convey("should return a parsed name usage", func() {
			n, err := species.ParsedName()
			So(err, ShouldBeNil)
			So(n.ScientificName, ShouldEqual, "Passer domesticus (Linnaeus, 1758)")
			So(n.Type, ShouldEqual, "SCIENTIFIC")
			So(n.AuthorsParsed, ShouldEqual, true)
		})

		Convey("should return a list of parent name usages", func() {
			r, err := species.Parents()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 6)
			So(r[0].ScientificName, ShouldEqual, "Animalia")
			So(r[0].KingdomKey, ShouldEqual, 1)
			So(r[0].AccordingTo, ShouldEqual, "The Catalogue of Life, 3rd January 2011")
		})

		Convey("should return a list of child name usages", func() {
			r, err := species.Children()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 12)
			So(r[0].DatasetKey, ShouldEqual, "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c")
			So(r[0].AccordingTo, ShouldEqual, "The Catalogue of Life, 3rd January 2011")
			So(r[0].CanonicalName, ShouldEqual, "Passer domesticus domesticus")
			So(r[0].Key, ShouldEqual, 5231191)
			So(r[0].TaxonID, ShouldEqual, "119127244")
			So(r[0].TaxonomicStatus, ShouldEqual, "ACCEPTED")
		})

		Convey("should return a list of related name usages", func() {
			r, err := species.Related()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 18)
			So(r[0].DatasetKey, ShouldEqual, "0938172b-2086-439c-a1dd-c21cb0109ed5")
			So(r[0].AccordingTo, ShouldEqual, "CoL2006/ITS")
			So(r[0].CanonicalName, ShouldEqual, "Passer domesticus")
			So(r[0].Key, ShouldEqual, 107872894)
			So(r[0].TaxonID, ShouldEqual, "10582565")
			So(r[0].TaxonomicStatus, ShouldEqual, "ACCEPTED")
			So(r[0].Species, ShouldEqual, "Passer domesticus")
		})

		Convey("should return a list of synonymical name usages", func() {
			r, err := species.Synonyms()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 0)
		})

		Convey("should return a list of distributions", func() {
			r, err := species.Distributions()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 3)
			So(r[0].LocationID, ShouldEqual, "CO:52")
			So(r[0].Locality, ShouldEqual, "Isla el Morro | Isla de Tumaco | Bocagrande")
			So(r[0].SourceTaxonKey, ShouldEqual, 114319110)
			So(r[0].Remarks, ShouldEqual, "Barrio La Florida (parque y sede Batallón de Infantería de Marina No. 2) Aeropuerto La Florida, antigua empresa maderera barrio El Morrito, faro Capitanía Puerto - Centro de Control y Contaminación del Pacífico (CCCP) playa turística peña El Morro y El Arco, barrio La Cordialidad parque y sede Corponariño; Parque Colón, parque Nariño, muelle Residencias - embarcadero a Bocagrande; Playas Hotel Las Lilianas")
		})

		Convey("should return a list of media", func() {
			r, err := species.Media()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 1)
			So(r[0].Type, ShouldEqual, "StillImage")
			So(r[0].Identifier, ShouldEqual, "http://upload.wikimedia.org/wikipedia/commons/d/d9/House_sparrowII.jpg")
			So(r[0].SourceTaxonKey, ShouldEqual, 100448484)
		})

		Convey("should return a list of references", func() {
			r, err := species.References()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 10)
			So(r[0].Citation, ShouldEqual, "(1996) database, NODC Taxonomic Code")
			So(r[0].Type, ShouldEqual, "taxon")
			So(r[0].Source, ShouldEqual, "Catalogue of Life")
			So(r[0].SourceTaxonKey, ShouldEqual, 110906793)
		})

		Convey("should return a list of vernacular names", func() {
			r, err := species.VernacularNames()
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 36)
			So(r[0].VernacularName, ShouldEqual, "English sparrow")
			So(r[0].Language, ShouldEqual, "")
			So(r[0].Source, ShouldEqual, "Global Invasive Species Database")
			So(r[0].SourceTaxonKey, ShouldEqual, 100220560)
		})

	})

}

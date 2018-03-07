package gbif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMatch(t *testing.T) {

	t.Parallel()

	Convey("should provide expected match results", t, func() {
		r, err := match(matchQuery{
			Verbose: true,
			Kingdom: "Plantae",
			Name:    "Oenante",
		})
		So(err, ShouldBeNil)
		So(r.ScientificName, ShouldEqual, "Oenanthe L.")
		So(r.Confidence, ShouldEqual, 85)
		So(r.MatchType, ShouldEqual, "FUZZY")
		So(len(r.Alternatives), ShouldEqual, 2)
	})

	//Convey("should provide expected search results", t, func() {
	//
	//	r, err := search(searchQuery{
	//		Q:    "Puma",
	//		Rank: []rank{rankGenus},
	//	})
	//	So(err, ShouldBeNil)
	//	So(len(r), ShouldEqual, 4)
	//	So(r[0].DatasetKey, ShouldEqual, "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c")
	//	So(r[0].ScientificName, ShouldEqual, "Puma Jardine, 1834")
	//
	//})

}

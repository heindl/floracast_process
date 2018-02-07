package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSearch(t *testing.T) {

	t.Parallel()

	Convey("should provide expected match results", t, func() {
		r, err := Match(MatchQuery{
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

	Convey("should provide expected search results", t, func() {

		r, err := Search(SearchQuery{
			Q:    "Puma",
			Rank: []Rank{RankGENUS},
		})
		So(err, ShouldBeNil)
		So(len(r), ShouldEqual, 56)
		So(r[0].DatasetKey, ShouldEqual, "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c")
		So(r[0].ScientificName, ShouldEqual, "Puma Jardine, 1834")

	})

}

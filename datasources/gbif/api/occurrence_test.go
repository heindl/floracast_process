package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrence(t *testing.T) {

	t.Parallel()

	Convey("should provide a list of occurrence search results", t, func() {
		r, err := Occurrences(OccurrenceSearchQuery{
			TaxonKey:           7205815,
			HasCoordinate:      true,
			HasGeospatialIssue: false,
			LastInterpreted:    "2000-01-01,2015-12-30",
		})
		So(err, ShouldBeNil)
		So(len(r), ShouldEqual, 1545)
		So(r[0].PublishingOrgKey, ShouldEqual, "28eb1a3f-1c15-4a95-931a-4af90ecb574d")
	})

}

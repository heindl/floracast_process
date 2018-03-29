package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOccurrence(t *testing.T) {

	t.Parallel()

	Convey("should fetch occurrences", t, func() {
		res, err := FetchOccurrences(context.Background(), datasources.TargetID("8229116"), nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 17)

	})

	Convey("should provide a list of occurrence search results", t, func() {
		r, err := fetchOccurrences(occurrenceSearchQuery{
			TaxonKey:           7205815,
			HasCoordinate:      true,
			HasGeospatialIssue: false,
			LastInterpreted:    "2000-01-01,2018-02-26",
		})
		So(err, ShouldBeNil)
		So(len(r), ShouldEqual, 607)
		So(r[0].PublishingOrgKey, ShouldEqual, "28eb1a3f-1c15-4a95-931a-4af90ecb574d")
	})

}

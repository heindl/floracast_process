package mushroomobserver

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"context"
	"bitbucket.org/heindl/processors/datasources"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("Should generate a list of NameUsageSources", t, func() {

		names := []string{
			"boletus edulis clavipes",
			"boletus edulis ochraceus",
			"boletus edulis grandedulis",
			"boletus edulis",
			"boletus elephantinus",
			"boletus edulis citrinus",
			"boletus slovenicus",
			"boletus edulis edulis",
			"boletus edulis arcticus",
			"boletus edulis arenarius",
			"boletus citrinus",
			"dictyopus edulis",
			"boletus edulis elephantinus",
			"leccinum elephantinum",
			"boletus edulis betulicola",
			"boletus edulis euedulis",
			"boletus venturii",
			"boletus edulis tardus",
			"boletus edulis quercus",
			"boletus quercicola",
			"boletus reticulatus albus",
			"boletus betulicola",
			"boletus edulis quercicola",
			"boletus edulis pseudopurpureus",
			"dictyopus edulis edulis",
			"boletus reticulatus citrinus",
			"boletus edulis praecox",
			"boletus solidus",
			"boletus edulis roseipes",
			"boletus edulis slovenicus",
			"suillus citrinus",
			"boletus edulis albus",
			"boletus edulis olivaceobrunneus",
			"boletus edulis subaereus",
			"boletus edulis laevipes",
			"boletus edulis trisporus",
			"leccinum edule",
			"tubiporus edulis",
			"boletus edulis piceicola",
			"boletus persoonii",
			"tylopilus porphyrosporus olivaceobrunneus",
			"boletus clavipes",
			"tubiporus edulis euedulis",
			"boletus edulis tuberosus",
			"boletus olivaceobrunneus",
			"tubiporus edulis edulis",
			"boletus edulis abietis",
			"boletus edulis communis",
		}

		res, err := MatchCanonicalNames(context.Background(), names...)
		So(err, ShouldBeNil)

		So(len(res), ShouldEqual, 4)

		requiredTargetIDs := datasources.TargetIDs{
			datasources.TargetID("16103"),
			datasources.TargetID("20594"),
			datasources.TargetID("16041"),
			datasources.TargetID("344"),
		}

		for _, src := range res {
			So(requiredTargetIDs.Contains(src.TargetID()), ShouldBeTrue)
		}

	})

	Convey("Should fetch MushroomObserver OccurrenceAggregation ", t, func() {

		res, err := FetchOccurrences(context.Background(), datasources.TargetID("16103"), nil)
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("20594"), nil)
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("16041"), nil)
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("344"), nil)
		So(err, ShouldBeNil)
		So(res.Count(), ShouldEqual, 16)

	})
}

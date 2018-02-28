package mushroomobserver

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaxonFetcher(t *testing.T) {

	Convey("Should generate a list of NameUsageSources", t, func() {

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

		res, err := FetchNameUsages(context.Background(), names, nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 4)

		for i, ex := range [][2]string{
			{"boletus edulis", "344"},
			{"boletus quercicola", "16103"},
			{"boletus betulicola", "20594"},
			{"boletus persoonii", "16041"},
		} {
			names, err := res[i].AllScientificNames()
			So(err, ShouldBeNil)
			So(names, ShouldContain, ex[0])
			hasSrc, err := res[i].HasSource(datasources.TypeMushroomObserver, datasources.TargetID(ex[1]))
			So(err, ShouldBeNil)
			So(hasSrc, ShouldBeTrue)
		}

	})

	SkipConvey("Should fetch MushroomObserver Occurrences", t, func() {

		res, err := FetchOccurrences(context.Background(), datasources.TargetID("16103"), nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("20594"), nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("16041"), nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 0)

		res, err = FetchOccurrences(context.Background(), datasources.TargetID("344"), nil)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 80)
		ids := []int{}
		// Check for duplicates
		for _, r := range res {
			ids = utils.AddIntToSet(ids, r.ID)
		}
		So(len(ids), ShouldEqual, 80)

	})
}

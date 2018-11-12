package gbif

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/utils"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNameUsageFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch gbif name usages from match strings", t, func() {
		cn := "morchella esculenta"
		taxa, err := FetchNameUsages(context.Background(), []string{cn}, nil)
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 1)
		So(taxa[0].CanonicalName().ScientificName(), ShouldEqual, cn)
		So(taxa[0].CanonicalName().Rank, ShouldEqual, "species")
		commonName, err := taxa[0].CommonName()
		So(err, ShouldBeNil)
		So(commonName, ShouldEqual, "morel")
		srcs, err := taxa[0].Sources()
		So(err, ShouldBeNil)
		So(len(srcs), ShouldEqual, 126)
		scientificNames, err := taxa[0].AllScientificNames()
		So(err, ShouldBeNil)
		So(len(scientificNames), ShouldEqual, 103)
		for _, sn := range []string{
			"morchella distans spathulata",
			"morchella conica conica",
			"morchella cylindrica",
			"morchella dunensis",
			"morchella esculenta dunensis",
			"morchella esculenta fulva",
			"morchella umbrina umbrina",
			"morchella vulgaris alba",
			"morchella esculenta alba",
			"morchella esculenta esculenta",
			"morchella esculenta conica",
			"morchella esculenta lutescens",
			"morchella esculenta rigida",
			"morchella esculenta roseostraminea",
			"morchella rotunda fulva",
			"morchella distans distans",
			"phallus tremelloides",
			"morilla esculenta",
			"morchella esculenta viridis",
			"phallus esculentus albus",
			"morchella esculenta albida",
			"morchella pubescens",
			"phalloboletus esculentus",
			"phallus esculentus cinereus",
			"morchella distans",
			"morchella conica cilicicae",
			"morchella abietina",
			"morchella esculenta rotunda",
			"morchella rotunda cinerea",
			"morchella esculenta violacea",
			"morchella dunensis dunensis",
			"morchella esculenta grisea",
			"morchella umbrina",
			"morilla conica",
			"morchella conica angusticeps",
			"morchella esculenta ovalis",
			"morchella tremelloides",
			"morchella vulgaris cinerascens",
			"morchella conica cylindrica",
			"morchella esculenta prunarii",
			"morchella umbrina macroalveola",
			"morchella viridis",
			"morilla tremelloides",
			"morchella conica elata",
			"morchella esculenta abietina",
			"morchella esculenta pubescens",
			"morchella rotunda minutella",
			"morchella conica pusilla",
			"morchella esculenta theobromichroa",
			"morchella rotunda rotunda",
			"morchella rotunda minutela",
			"morchella vulgaris parvula",
			"morchella dunensis sterile",
			"morchella esculenta umbrinoides",
			"morchella esculenta rubroris",
			"phallus esculentus fuscus",
			"morchella rotunda alba",
			"morchella esculenta aurantiaca",
			"morchella esculenta brunnea",
			"morchella esculenta umbrina",
			"morchella ovalis pallida",
			"phallus esculentus rotundus",
			"morchella conica",
			"morchella conica crassa",
			"morchella rigida",
			"helvella esculenta",
			"morellus esculentus",
			"phallus esculentus",
			"morchella lutescens",
			"morchella conica rigida",
			"morchella rotunda crassipes",
			"morchella vulgaris",
			"morchella conica pygmaea",
			"morchella esculenta sterilis",
			"morchella esculenta mahoniae",
			"morchella dunensis sterilis",
			"morchella rotunda",
			"morchella rotunda esculenta",
			"morchella vulgaris vulgaris",
			"morchella esculenta stipitata",
			"morchella conica ceracea",
			"morchella distans longissima",
			"morchella rotunda minutula",
			"morchella rotunda pallida",
			"morchella esculenta atrotomentosa",
			"morchella vulgaris albida",
			"morchella conica distans",
			"morchella rotunda olivacea",
			"morchella rotunda pubescens",
			"morchella atrotomentosa",
			"morchella conica nigra",
			"morchella conica violeipes",
			"morchella prunarii",
			"morchella conica metheformis",
			"morchella conica flexuosa",
			"morchella esculenta corrugata",
			"morchella esculenta vulgaris",
			"morchella rotunda rigida",
			"phallus esculentus esculentus",
			"morchella vulgaris tremelloides",
			"morchella conica serotina",
			"morchella conica meandriformis",
			"morchella esculenta",
		} {
			So(utils.ContainsString(scientificNames, sn), ShouldBeTrue)
		}

	})

	Convey("should fetch gbif name usages from keys", t, func() {

		names := []string{
			"morchella brunnea",
			"morchella snyderi",
			"morchella importuna",
			"morchella diminutiva",
			"morchella frustrata",
			"morchella esculentoides",
			"morchella americana",
			"morchella tridentina",
			"morchella costata",
			"morchella guatemalensis",
			"morchella elata",
			"morchella deliciosa",
			"morchella esculenta",
			"morchella tomentosa",
			"morchella eximia",
			"morchella populiphila",
			"morchella pragensis",
			"morchella vulgaris",
			"morchella australiana",
			"morchella septimelata",
			"morchella sextelata",
			"morchella crassipes",
			"morchella conica",
			"morchella angusticeps",
			"morchella punctipes",
			"morchella umbrina",
			"morchella prava",
			"morchella rufobrunnea",
		}

		ids := datasources.TargetIDs{"2594601",
			"3495681",
			"2594630",
			"3495624",
			"2594626",
			"2594620",
			"2594602",
			"3495605",
			"7258337",
			"8395134",
			"6019620",
			"2594615",
			"8029683",
			"7468887",
			"7561838",
			"2594612",
			"2594617",
			"2594624",
			"3495647",
			"2594604",
			"3495648"}

		taxa, err := FetchNameUsages(context.Background(), names, ids)
		So(err, ShouldBeNil)
		So(len(taxa), ShouldEqual, 27)

	})
}

package gbif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch gbif name usages from match strings", t, func() {
		taxa, err := FetchNamesUsages(context.Background(),  []string{"Morchella esculenta"}, nil)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(taxa))
	})

	Convey("should fetch gbif name usages from keys", t, func() {
		m := map[TaxonID][]string{
			TaxonID(2594620): []string{
				"Morchella deliciosa",
			},
			TaxonID(2594615): []string{
				"Morchella vulgaris",
			},
			TaxonID(3495647): []string{
				"Morchella punctipes",
			},
			TaxonID(3495605): []string{
				"Morchella tomentosa",
			},
			TaxonID(2594630): []string{
				"Morchella costata",
			},
			TaxonID(6019620): []string{
				"Morchella pragensis",
			},
			TaxonID(3495681): []string{
				"Morchella tridentina",
			},
			TaxonID(2594604): []string{
				"Morchella umbrina",
			},
			TaxonID(2594624): []string{
				"Morchella angusticeps",
			},
			TaxonID(3495624): []string{
				"Morchella guatemalensis",
			},
			TaxonID(2594612): []string{
				"Morchella crassipes",
			},
			TaxonID(8395134): []string{
				"Morchella populiphila",
			},
			TaxonID(2594601): []string{
				"Morchella esculentoides",
				"Morchella diminutiva",
				"Morchella importuna",
				"Morchella americana",
				"Morchella snyderi",
				"Morchella frustrata",
				"Morchella brunnea",
			},
			TaxonID(7258337): []string{
				"Morchella eximia",
			},
			TaxonID(2594617): []string{
				"Morchella conica",
			},
			TaxonID(2594626): []string{
				"Morchella elata",
			},
			TaxonID(7468887): []string{
				"Morchella septimelata",
			},
			TaxonID(2594602): []string{
				"Morchella esculenta",
			},
			TaxonID(8029683): []string{
				"Morchella australiana",
			},
			TaxonID(3495648): []string{
				"Morchella rufobrunnea",
			},
			TaxonID(7561838): []string{
				"Morchella sextelata",
			},
		}

		taxa, err := FetchNamesUsages(context.Background(),  nil, m)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(taxa))
	})
}

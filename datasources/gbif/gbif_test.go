package gbif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"fmt"
	"bitbucket.org/heindl/process/utils"
	"bitbucket.org/heindl/process/datasources"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch occurrences", t, func() {
		res, err := FetchOccurrences(context.Background(), datasources.TargetID("8229116"), nil)
		So(err, ShouldBeNil)
		fmt.Println("Result", len(res))

	})

	SkipConvey("should fetch gbif name usages from match strings", t, func() {
		taxa, err := FetchNamesUsages(context.Background(),  []string{"Morchella esculenta"}, nil)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(taxa))
	})

	Convey("should fetch photos", t, func() {
		photos, err := FetchPhotos(context.Background(),  datasources.TargetID("2594602"))
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(photos))
	})

	SkipConvey("should fetch gbif name usages from keys", t, func() {

		names := []string {
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

		taxa, err := FetchNamesUsages(context.Background(), names, ids)
		So(err, ShouldBeNil)

		fmt.Println("Result", len(taxa))

	})
}

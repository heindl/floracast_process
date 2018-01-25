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
		m := map[int][]string{
			5249503: []string{
				"Cantharellus lutescens",
			},
			7701489: []string{
				"Cantharellus phasmatis",
			},
			5184831: []string{
				"Cantharellus noumeae",
			},
			7509721: []string{
				"Cantharellus solidus",
			},
			8186696: []string{
				"Cantharellus nigrescens",
			},
			5249568: []string{
				"Cantharellus camphoratus",
			},
			5249528: []string{
				"Cantharellus cinereus",
			},
			5249595: []string{
				"Cantharellus atrolilacinus",
			},
			5467112: []string{
				"Cantharellus congolensis",
			},
			5184832: []string{
				"Cantharellus jebbi",
			},
			5249489: []string{
				"Cantharellus wellingtonensis",
			},
			7241872: []string{
				"Cantharellus subcibarius",
			},
			7241861: []string{
				"Cantharellus lateritius",
			},
			5466942: []string{
				"Cantharellus tabernensis",
			},
			5466954: []string{
				"Cantharellus formosus",
			},
			5249502: []string{
				"Cantharellus elsae",
			},
			5184830: []string{
				"Cantharellus doederleini",
			},
			5249497: []string{
				"Cantharellus cascadensis",
			},
			5249504: []string{
				"Cantharellus cibarius",
			},
			5249564: []string{
				"Cantharellus minor",
			},
			8269489: []string{
				"Cantharellus lewisii",
			},
			8345326: []string{
				"Cantharellus variabilicolor",
			},
			5466960: []string{
				"Cantharellus subalbidus",
			},
			5249498: []string{
				"Cantharellus subpruinosus",
			},
			5249555: []string{
				"Cantharellus amethysteus",
			},
			5466926: []string{
				"Cantharellus goossensiae",
			},
			5249578: []string{
				"Cantharellus insignis",
			},
			7488096: []string{
				"Cantharellus tricolor",
			},
			5467036: []string{
				"Cantharellus appalachiensis",
			},
			7994006: []string{
				"Cantharellus flavus",
			},
			5467048: []string{
				"Cantharellus variabilis",
			},
			5249494: []string{
				"Cantharellus ferruginascens",
			},
			5249601: []string{
				"Cantharellus persicinus",
			},
			5249602: []string{
				"Cantharellus attenuatus",
			},
			7241881: []string{
				"Cantharellus odoratus odoratus",
			},
			5249583: []string{
				"Cantharellus ignicolor",
			},
			5249533: []string{
				"Cantharellus subincarnatus",
			},
			5249565: []string{
				"Cantharellus defibulatus",
			},
			5249496: []string{
				"Cantharellus pallens",
			},
			8015230: []string{
				"Cantharellus spectaculus",
			},
			5249599: []string{
				"Cantharellus cinnabarinus",
			},
			7443622: []string{
				"Cantharellus roseocanus",
			},
		}

		taxa, err := FetchNamesUsages(context.Background(),  nil, m)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(taxa))
	})
}

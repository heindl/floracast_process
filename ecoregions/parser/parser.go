package main

import (
	"bitbucket.org/heindl/processors/terra"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"os"
)

//
//type CachedEcoRegion struct {
//	Realm   string `json:"REALM"`
//	EcoName string `json:"ECO_NAME"` // Ecoregion Name
//	EcoCode string `json:"eco_code"` // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
//	EcoNum  int64  `json:"ECO_NUM"`  // A unique number for each ecoregion within each biome nested within each realm.
//	EcoID   int64  `json:"ECO_ID"`   // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
//	//GlobalStatus int64 `json:"GBL_STAT"`
//	Biome int64 `json:"BIOME"`
//	//Area float64 `json:"AREA"`
//	//AreaKm2 float64 `json:"area_km2"`
//	MultiPolygon terra.MultiPolygon
//}
//
//const cachedEcoRegionString = `
//	type CachedEcoRegion struct {
//		Realm  string
//		EcoName string
//		EcoCode string
//		EcoNum int64
//		EcoID int64
//		Biome int64
//		MultiPolygon terra.MultiPolygon
//	}
//`

func main() {

	ecoRegionFile := flag.String("ecoregions", "/Users/m/Downloads/wwf_terr_ecos_oRn.json", "Path to ecoregions json file.")

	flag.Parse()

	// File can be downloaded here: https://worldmap.harvard.edu/data/geonode:wwf_terr_ecos_oRn
	// http://worldmap.harvard.edu/download/wfs/697/json?outputFormat=json&service=WFS&request=GetFeature&format_options=charset%3AUTF-8&typename=geonode%3Awwf_terr_ecos_oRn&version=1.0.0

	fc, err := terra.ReadFeatureCollectionFromGeoJSONFile(*ecoRegionFile, func(b []byte) bool {
		r, err := jsonparser.GetString(b, "REALM")
		if err != nil {
			panic(err)
		}
		return r != "NA"
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("After initial retrieval", fc.Count())

	fc_id_grouped := fc.GroupByProperties("ECO_ID")

	fc_id_condensed, err := fc_id_grouped.Condense(func(a, b []byte) []byte {
		// Note that we are only storing categories, and those should all be the same as we're grouping by eco_id.
		return b
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("After Group and Condense", fc_id_condensed.Count())

	b, err := fc_id_condensed.GeoJSON()
	if err != nil {
		panic(err)
	}

	content := []byte("package ecoregions\nconst ecoregions_geojson=`")
	content = append(content, b...)
	content = append(content, []byte("`")...)

	if err := ioutil.WriteFile("./geojson.go", content, os.ModePerm); err != nil {
		panic(err)
	}

	//
	//fmt.Println(fmt.Sprintf("Feature Collection Area: %.6f", fc_id_condensed.Area()))
	//
	//fname := fmt.Sprintf("../cache/Cache.go")
	//main_file, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	//if err != nil {
	//	panic(err)
	//}
	//defer main_file.Close()
	//header := fmt.Sprintf(`
	//	package cache
	//
	//	%s
	//
	//	var EcoRegionCache = []CachedEcoRegion{}
	//	`, cachedEcoRegionString)
	//header = strings.TrimSpace(header)
	//if _, err = main_file.WriteString(header); err != nil {
	//	panic(err)
	//}
	//
	//for _, feature := range fc_id_condensed.Features() {
	//
	//	cr := CachedEcoRegion{}
	//	if err := feature.GetProperties(&cr); err != nil {
	//		panic(err)
	//	}
	//
	//	cr.MultiPolygon = feature.MultiPolygon()
	//
	//	fmt.Println("Writing file", cr.EcoID)
	//	func() {
	//		fname := fmt.Sprintf("../cache/EcoID_%d.go", cr.EcoID)
	//		nf, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	//		if err != nil {
	//			panic(err)
	//		}
	//		defer nf.Close()
	//
	//		page := fmt.Sprintf(`
	//			package cache
	//
	//			func init() {
	//				EcoRegionCache = append(EcoRegionCache, %#v)
	//			}
	//
	//			`, cr)
	//		page = strings.Replace(page, "main.", "", -1)
	//		if _, err = nf.WriteString(page); err != nil {
	//			panic(err)
	//		}
	//	}()
	//
	//}

}

//type Parser struct {
//	sync.Mutex
//	File       *os.File
//	FileMap    map[int64]*CachedEcoRegion
//	AllRegions []CachedEcoRegion
//}
//
//func (Ω *Parser) PrintDefinitions() {
//
//	idNames := map[int64]string{}
//	codeNames := map[string]string{}
//	biomes := map[int64]int{}
//
//	for _, r := range Ω.AllRegions {
//		idNames[r.EcoID] = r.EcoName
//		codeNames[r.EcoCode] = r.EcoName
//		biomes[r.Biome] += 1
//	}
//
//	fmt.Println("var EcoIDDefinitions = map[EcoID]string{")
//	for k, v := range idNames {
//		fmt.Println(fmt.Sprintf(`EcoID(%d): "%s",`, k, v))
//	}
//	fmt.Println("}")
//
//	fmt.Println("var EcoCodeDefinitions = map[EcoCode]string{")
//	for k, v := range codeNames {
//		fmt.Println(fmt.Sprintf(`EcoCode("%s"): "%s",`, k, v))
//	}
//	fmt.Println("}")
//
//	fmt.Println(utils.JsonOrSpew(biomes))
//
//	return
//}

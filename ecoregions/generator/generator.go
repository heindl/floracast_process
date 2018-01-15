package main

import (
	"github.com/saleswise/errors/errors"
	"bitbucket.org/heindl/taxa/utils"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"bitbucket.org/heindl/taxa/terra"
)

type CachedEcoRegion struct {
	Realm   string `json:"REALM"`
	EcoName string `json:"ECO_NAME"` // Ecoregion Name
	EcoCode string `json:"eco_code"` // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum  int64  `json:"ECO_NUM"`  // A unique number for each ecoregion within each biome nested within each realm.
	EcoID   int64  `json:"ECO_ID"`   // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	//GlobalStatus int64 `json:"GBL_STAT"`
	Biome int64 `json:"BIOME"`
	//Area float64 `json:"AREA"`
	//AreaKm2 float64 `json:"area_km2"`
	MultiPolygon terra.MultiPolygon
}

const cachedEcoRegionString = `
	type CachedEcoRegion struct {
		Realm  string
		EcoName string
		EcoCode string
		EcoNum int64
		EcoID int64
		Biome int64
		Geometries []string
	}
`

func main() {

	ecoRegionFile := flag.String("ecoregions", "/Users/m/Downloads/wwf_terr_ecos_oRn.json", "Path to ecoregions json file.")

	flag.Parse()

	p := Parser{
		FileMap:    map[int64]*CachedEcoRegion{},
		AllRegions: []CachedEcoRegion{},
	}

	// File can be downloaded here: https://worldmap.harvard.edu/data/geonode:wwf_terr_ecos_oRn
	// http://worldmap.harvard.edu/download/wfs/697/json?outputFormat=json&service=WFS&request=GetFeature&format_options=charset%3AUTF-8&typename=geonode%3Awwf_terr_ecos_oRn&version=1.0.0

	if err := terra.ReadGeoJSONFeatureCollectionFile(*ecoRegionFile, p.Parse); err != nil {
		panic(err)
	}

	//p.PrintDefinitions()
	//return

	fname := fmt.Sprintf("../generated_cache/Cache.go")
	main_file, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer main_file.Close()
	header := fmt.Sprintf(`
		package cache

		%s

		var EcoRegionCache = []CachedEcoRegion{}
		`, cachedEcoRegionString)
	header = strings.TrimSpace(header)
	if _, err = main_file.WriteString(header); err != nil {
		panic(err)
	}

	for k, region := range p.FileMap {
		func() {
			fname := fmt.Sprintf("../generated_cache/EcoID_%d.go", k)
			nf, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				panic(err)
			}
			defer nf.Close()

			page := fmt.Sprintf(`
				package cache
		
				func init() {
					EcoRegionCache = append(EcoRegionCache, %#v)
				}
		
				`, *region)
			page = strings.Replace(page, "main.", "", -1)
			if _, err = nf.WriteString(page); err != nil {
				panic(err)
			}
		}()

	}

}

type Parser struct {
	sync.Mutex
	File       *os.File
	FileMap    map[int64]*CachedEcoRegion
	AllRegions []CachedEcoRegion
}

func (Ω *Parser) Parse(encoded_properties []byte, multipolygon terra.MultiPolygon) error {


	o := CachedEcoRegion{}
	if err := json.Unmarshal(encoded_properties, &o); err != nil {
		return errors.Wrap(err, "could not unmarshal ecoregion")
	}

	if o.Realm != "NA" {
		return nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	if _, ok := Ω.FileMap[o.EcoID]; !ok {
		Ω.FileMap[o.EcoID] = &o
	}
	Ω.FileMap[o.EcoID].MultiPolygon = Ω.FileMap[o.EcoID].MultiPolygon.PushMultiPolygon(multipolygon)

	Ω.AllRegions = append(Ω.AllRegions, o)

	return nil
}

func (Ω *Parser) PrintDefinitions() {

	idNames := map[int64]string{}
	codeNames := map[string]string{}
	biomes := map[int64]int{}

	for _, r := range Ω.AllRegions {
		idNames[r.EcoID] = r.EcoName
		codeNames[r.EcoCode] = r.EcoName
		biomes[r.Biome] += 1
	}

	fmt.Println("var EcoIDDefinitions = map[EcoID]string{")
	for k, v := range idNames {
		fmt.Println(fmt.Sprintf(`EcoID(%d): "%s",`, k, v))
	}
	fmt.Println("}")

	fmt.Println("var EcoCodeDefinitions = map[EcoCode]string{")
	for k, v := range codeNames {
		fmt.Println(fmt.Sprintf(`EcoCode("%s"): "%s",`, k, v))
	}
	fmt.Println("}")

	fmt.Println(utils.JsonOrSpew(biomes))

	return
}

package main

import (
	"bufio"
	"github.com/buger/jsonparser"
	"github.com/saleswise/errors/errors"
	//geojson "github.com/tidwall/tile38/geojson"
	"io/ioutil"
	"os"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"github.com/tidwall/tile38/geojson"
	"sync"
)

type CachedEcoRegion struct {
	Realm  string `json:"REALM"`
	//EcoName string `json:"ECO_NAME"` // Ecoregion Name
	EcoCode string `json:"eco_code"` // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum int64 `json:"ECO_NUM"` // A unique number for each ecoregion within each biome nested within each realm.
	EcoID int64 `json:"ECO_ID"` // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	//GlobalStatus int64 `json:"GBL_STAT"`
	Biome int64 `json:"BIOME"`
	//Area float64 `json:"AREA"`
	//AreaKm2 float64 `json:"area_km2"`
	GeoHashes []string
	GeoJsonString string
}
const cachedEcoRegionString = `
	type CachedEcoRegion struct {
		Realm  string
		//EcoName string
		EcoCode string
		EcoNum int64
		EcoID int64
		//GlobalStatus int64
		Biome int64
		//Area float64
		//AreaKm2 float64
		GeoHashes []string
		GeoJsonString string
	}
`



func main() {

	ecoRegionFile := flag.String("ecoregions", "/Users/m/Downloads/wwf_terr_ecos_oRn.json", "Path to ecoregions json file.")

	flag.Parse()

	// File can be downloaded here: https://worldmap.harvard.edu/data/geonode:wwf_terr_ecos_oRn
	// http://worldmap.harvard.edu/download/wfs/697/json?outputFormat=json&service=WFS&request=GetFeature&format_options=charset%3AUTF-8&typename=geonode%3Awwf_terr_ecos_oRn&version=1.0.0

	ecoFile, err := os.Open(*ecoRegionFile)
	if err != nil {
		panic(err)
	}
	defer ecoFile.Close()

	b, err := ioutil.ReadAll(bufio.NewReader(ecoFile))
	if err != nil {
		panic(err)
	}

	p := Parser{
		FileMap: map[int64][]string{},
	}
	if _, err := jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
		if err != nil {
			panic(err)
		}
		if err := p.Parse(value); err != nil {
			panic(err)
		}
	}, "features"); err != nil {
		panic(errors.Wrap(err, "could not parse features array"))
	}

	fname := fmt.Sprintf("../generated_cache/generated_cache.go")
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


	for k, rows := range p.FileMap {
		func() {
			fname := fmt.Sprintf("../generated_cache/generated_regions_%d.go", k)
			nf, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				panic(err)
			}
			defer nf.Close()

			header := `
		package cache

		func init() {
			EcoRegionCache = append(EcoRegionCache, []CachedEcoRegion{
		`
			if _, err = nf.WriteString(header); err != nil {
				panic(err)
			}
			for _, r := range rows {
				if _, err = nf.WriteString(r + "\n"); err != nil {
					panic(err)
				}
			}
			footer := `
				}...)
			}
			`
			if _, err = nf.WriteString(footer); err != nil {
				panic(err)
			}
		}()

	}

}

type Parser struct {
	sync.Mutex
	File *os.File
	FileMap map[int64][]string
}

func (Ω *Parser) Parse(value []byte) error {

		b, _, _, err := jsonparser.Get(value, "properties")
		if err != nil {
			return errors.Wrap(err, "could not get properties")
		}

		o := CachedEcoRegion{}
		if err := json.Unmarshal(b, &o); err != nil {
			return errors.Wrap(err, "could not unmarshal ecoregion")
		}

		if o.Realm != "NT" {
			return nil
		}

		geometry, _, _, err := jsonparser.Get(value, "geometry")
		if err != nil {
			return errors.Wrap(err, "could not get geometry")
		}

		obj, err := geojson.ObjectJSON(string(geometry))
		if err != nil {
			return errors.Wrap(err, "could not parse geojson object")
		}

		if ok := geojson.New2DPoint(38.6270025, -90.1994042).Within(obj); ok {
			fmt.Println("is within")
		}

		for i := 1; i <= 10; i++ {
			hash, err := obj.Geohash(i)
			if err != nil {
				return errors.Wrap(err, "could not get geohash")
			}
			o.GeoHashes = append(o.GeoHashes, hash)
		}

		o.GeoJsonString = obj.String()

		s := strings.Replace(fmt.Sprintf("%#v,", o), "main.CachedEcoRegion", "", -1)

		Ω.Lock()
		defer Ω.Unlock()

		if _, ok := Ω.FileMap[o.Biome]; !ok {
			Ω.FileMap[o.Biome] = []string{}
		}

		Ω.FileMap[o.Biome] = append(Ω.FileMap[o.Biome], s)

		return nil
}
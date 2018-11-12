package main

import (
	"github.com/heindl/floracast_process/terra/geo"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {

	ecoRegionFile := flag.String("ecoregions", "/Users/m/Downloads/wwf_terr_ecos_oRn.json", "Path to ecoregions json file.")

	flag.Parse()

	// File can be downloaded here: https://worldmap.harvard.edu/data/geonode:wwf_terr_ecos_oRn
	// http://worldmap.harvard.edu/download/wfs/697/json?outputFormat=json&service=WFS&request=GetFeature&format_options=charset%3AUTF-8&typename=geonode%3Awwf_terr_ecos_oRn&version=1.0.0

	fc, err := geo.ReadFeatureCollectionFromGeoJSONFile(*ecoRegionFile, func(b []byte) bool {
		id, err := jsonparser.GetInt(b, "ECO_ID")
		if err != nil {
			panic(err)
		}

		sID := strings.TrimSpace(strconv.Itoa(int(id)))
		if strings.Contains(sID, "-") {
			return true
		}

		if sID[:1] != "5" && sID[:1] != "6" {
			return true
		}

		return false
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%d Records After Initial Parse & Filter", fc.Count()))

	fcIDGrouped, err := fc.GroupByProperties("ECO_ID")
	if err != nil {
		panic(err)
	}

	fcIDCondensed, err := fcIDGrouped.Condense(func(a, b []byte) []byte {
		// Note that we are only storing categories, and those should all be the same as we're grouping by eco_id.
		return b
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%d Records After Group & Condense", fcIDCondensed.Count()))

	b, err := fcIDCondensed.GeoJSON()
	if err != nil {
		panic(err)
	}

	content := []byte("package cache\nconst ecoregionsGeoJSON=`")
	content = append(content, b...)
	content = append(content, []byte("`")...)

	if err := ioutil.WriteFile("./geojson.go", content, os.ModePerm); err != nil {
		panic(err)
	}

}

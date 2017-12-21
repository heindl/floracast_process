package main

import (
	"flag"
	"io/ioutil"
	"github.com/paulmach/go.geojson"
	"fmt"
	"strconv"
	"encoding/json"
	"path"
)

func main() {
	in := flag.String("in", "", "Input json file.")
	out := flag.String("out", "", "Combined json directory.")

	flag.Parse()

	if *in == "" || *out == "" {
		panic("input file and output directory required.")
	}

	b, err := ioutil.ReadFile(*in)
	if err != nil {
		panic(err)
	}

	inputFeatureCollection, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		panic(err)
	}

	outputFeatureCollections := map[string]*geojson.FeatureCollection{}

	zeros := 0
	for _, f := range inputFeatureCollection.Features{
		if v, ok := f.Properties["WDPA_Cd"]; ok {
			fid := v.(float64)
			if fid == 0 {
				zeros++
				continue
			}
			id := strconv.FormatFloat(fid, 'f', 0, 64)
			if _, ok := outputFeatureCollections[id]; !ok {
				outputFeatureCollections[id] = &geojson.FeatureCollection{}
			}
			outputFeatureCollections[id].AddFeature(f)
		}
	}

	fmt.Println("ZERO COUNT", zeros)
	fmt.Println("TOTAL AREAS", len(outputFeatureCollections))

	for k, ofc := range outputFeatureCollections {
		b, err := json.Marshal(ofc)
		if err != nil {
			panic(err)
		}
		fname := k + ".json"
		if err := ioutil.WriteFile(path.Join(*out, fname), b, 0700); err != nil {
			panic(err)
		}
	}


}
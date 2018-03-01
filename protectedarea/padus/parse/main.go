package main

import (
	"flag"
)

var flagsToFilter = []string{
	"golf", "soccer", "recreation", "athletic", "softball",
	"baseball", "horse", "arts", "gym", "cemetery", "museum",
	"community center", "sports", "city high", "high school", "tennis", "pavilion",
	"skate park", "unknown", "elementary", "library",
}

const maxCentroidDistance = 20.0
const clusterDecimationKm = 10.0

func main() {
	in := flag.String("in", "/tmp/gap_analysis/ID/state.geojson", "Input json file")
	out := flag.String("out", "/tmp/gap_analysis/ID/areas", "Combined json directory")

	flag.Parse()

	process, err := NewProcessor(*in, *out)
	if err != nil {
		panic(err)
	}

	collections, _, err := process.ProcessFeatureCollections()
	if err != nil {
		panic(err)
	}

	if err := process.WriteCollections(collections); err != nil {
		panic(err)
	}

}

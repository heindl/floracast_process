package main

import (
	"flag"
	"fmt"
	"bitbucket.org/heindl/process/store"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/process/occurrences"
)

func main() {
	number := flag.Float64("number", 0, "Total number of random occurrences to generate")
	flag.Parse()

	if *number == 0 {
		panic("Please include the number of random occurrences to generate. Will be increased to at least a multiple of four from the number of grid divisions.")
	}

	aggr, err := occurrences.GenerateRandomOccurrences(*number)
	if err != nil {
		panic(err)
	}

	if aggr == nil {
		return
	}

	fmt.Println("Count", aggr.Count())

	geojson, err := aggr.GeoJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(geojson))

	//return

	cxt := context.Background()

	florastore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	if err := occurrences.ClearRandomOccurrences(cxt, florastore); err != nil {
		panic(err)
	}

	if err := aggr.Upload(cxt, florastore); err != nil {
		panic(err)
	}

}

package main

import (
	"github.com/heindl/floracast_process/store"
	"context"
	"flag"
)

func main() {

	// /tmp/gap_analysis/OR/areas

	geojsonPath := flag.String("path", "", "Path to geojson files to search recursively.")

	flag.Parse()

	if *geojsonPath == "" {
		panic("A geojson directory must be specified.")
	}

	areas, err := ParseProtectedAreaDirectory(*geojsonPath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		panic(err)
	}

	if _, err := areas.Upload(ctx, floraStore); err != nil {
		panic(err)
	}
}

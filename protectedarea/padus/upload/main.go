package main

import (
	"bitbucket.org/heindl/process/store"
	"flag"
	"golang.org/x/net/context"
)

func main() {

	geojsonPath := flag.String("geojson", "/tmp/gap_analysis", "Path to geojson files to search recursively.")

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

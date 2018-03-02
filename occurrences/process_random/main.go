package main

import (
	"bitbucket.org/heindl/process/occurrences"
	"bitbucket.org/heindl/process/store"
	"context"
	"flag"
)

func main() {
	batches := flag.Int("batches", 0, "Random point batches")
	cellLevel := flag.Int("level", 0, "S2 cell level")
	flag.Parse()

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		panic(err)
	}

	aggr, err := occurrences.GenerateRandomOccurrences(*cellLevel, *batches)
	if err != nil {
		panic(err)
	}

	if err := aggr.Upload(context.Background(), floraStore); err != nil {
		panic(err)
	}

}

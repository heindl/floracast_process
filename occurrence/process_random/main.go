package main

import (
	"bitbucket.org/heindl/process/occurrence"
	"bitbucket.org/heindl/process/store"
	"context"
	"flag"
)

// Level 4: 216
// Level 3: 64

func main() {
	batches := flag.Int("batches", 0, "Random point batches")
	cellLevel := flag.Int("level", 0, "S2 cell level")
	flag.Parse()

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		panic(err)
	}

	if err = occurrence.ClearRandomPoints(ctx, floraStore); err != nil {
		panic(err)
	}

	aggr, err := occurrence.GenerateRandomOccurrences(*cellLevel, *batches)
	if err != nil {
		panic(err)
	}

	if err = aggr.Upload(context.Background(), floraStore); err != nil {
		panic(err)
	}

}

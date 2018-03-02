package main

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/occurrence"
	"bitbucket.org/heindl/process/store"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		panic(err)
	}

	nameUsageIterator := nameusage.FetchAll(ctx, floraStore)

	for {
		usage, err := nameUsageIterator.Next()
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		aggr, err := occurrence.FetchOccurrences(ctx, usage, true)
		if err != nil {
			panic(err)
		}

		if err := aggr.Upload(ctx, floraStore); err != nil {
			panic(err)
		}

		// Update usage with new fetch times.
		if _, err := usage.Upload(ctx, floraStore); err != nil {
			panic(err)
		}
	}

}

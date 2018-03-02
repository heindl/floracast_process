package main

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/occurrences"
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
		fetchSource(ctx, floraStore, usage)
	}

}

func fetchSource(ctx context.Context, floraStore store.FloraStore, usage nameusage.NameUsage) {
	srcs, err := usage.Sources()
	if err != nil {
		panic(err)
	}
	for _, src := range srcs {
		if src.OccurrenceCount() > 0 {

			srcType, err := src.SourceType()
			if err != nil {
				panic(err)
			}
			targetID, err := src.TargetID()
			if err != nil {
				panic(err)
			}
			aggr, err := occurrences.FetchOccurrences(ctx, srcType, targetID, src.LastFetchedAt())
			if err != nil {
				panic(err)
			}
			if err := aggr.Upload(ctx, floraStore); err != nil {
				panic(err)
			}
		}
	}
}

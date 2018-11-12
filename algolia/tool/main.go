package main

import (
	"github.com/heindl/floracast_process/algolia"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"flag"
	"golang.org/x/net/context"
)

func main() {

	usageStrPtr := flag.String("usage", "", "NameUsageId to materialize NameUsage for")
	flag.Parse()

	nameUsageID := nameusage.ID(*usageStrPtr)
	if !nameUsageID.Valid() {
		panic("Invalid NameUsageID")
	}

	cxt := context.Background()

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	if err := algolia.IndexNameUsage(cxt, floraStore, nameUsageID); err != nil {
		panic(err)
	}

}

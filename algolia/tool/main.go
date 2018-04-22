package main

import (
	"bitbucket.org/heindl/process/algolia"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
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

package main

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/sourcefetchers"
	"bitbucket.org/heindl/process/nameusage/aggregate"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/occurrence"
	"bitbucket.org/heindl/process/store"
	"context"
	"flag"
	"github.com/golang/glog"
	"strconv"
	"strings"
)

const minimumOccurrenceCount = 100

func main() {

	ids := flag.String("inaturalists", "", "Comma separated INaturalist Taxon integer ids.")

	flag.Parse()

	if *ids == "" {
		return
	}

	intIDs := []int{}
	for _, id := range strings.Split(*ids, ",") {
		i, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		intIDs = append(intIDs, i)
	}

	cxt := context.Background()

	nameUsageAggr, err := aggregateInitialNameUsages(cxt, intIDs...)
	if err != nil {
		panic(err)
	}

	occurrenceAggr, err := fetchOccurrences(cxt, nameUsageAggr)
	if err != nil {
		panic(err)
	}

	fc, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	if err := nameUsageAggr.Upload(cxt, fc); err != nil {
		panic(err)
	}

	if err := occurrenceAggr.Upload(cxt, fc); err != nil {
		panic(err)
	}

}

// This function is intended to be a process of building a store taxon representation based on several
// sorted that accumulate taxonomic information. These names and synonyms can then be used to fetch occurrences
// from various sources using only the list of names.

// Note that this function aggressively combines synonyms at this point. If any of the three sources consider it a synonym,
// we accept it in order error on the side of more occurrences.

// Initially don't worry about existing values in the database. This function runs alone each time, connects to multiple
// datasources, and then deletes those in the database that match.

// The main goal here is occurrences. In order to do that we have to know which names are for the same species,
// so that sources with different occurrences for different species can be combined.

// TODO: One day convert all of this to a graph database. Which could even incorporate DNA records.
// In this scenario, all occurrences would need photos, and the photos would be used instead of names to connect them
// to phenotypical dna.

// TODO: Include a step to verify the ids for occurrences haven't changed.
// Could actually be done in the occurrence fetch step: https://www.inaturalist.org/taxon_changes

func aggregateInitialNameUsages(cxt context.Context, inaturalistTaxonIDs ...int) (*aggregate.Aggregate, error) {

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. TotalOccurrenceCount are still stored within the synonym.

	targetIDs, err := datasources.NewDataSourceTargetIDFromInts(datasources.TypeINaturalist, inaturalistTaxonIDs...)
	if err != nil {
		return nil, err
	}

	glog.Info("Aggregating a list of NameUsages from INaturalist TaxonIDs: ", targetIDs)

	snowball := aggregate.Aggregate{}

	// Order is obviously extremely important here.
	for _, srcType := range []datasources.SourceType{datasources.TypeINaturalist, datasources.TypeGBIF, datasources.TypeNatureServe} {

		ids, err := snowball.TargetIDs(srcType)
		if err != nil {
			return nil, err
		}
		if srcType == datasources.TypeINaturalist {
			ids = targetIDs
		}

		sciNames, err := snowball.ScientificNames()
		if err != nil {
			return nil, err
		}

		usages, err := sourcefetchers.FetchNameUsages(cxt, srcType, sciNames, ids)
		if err != nil {
			return nil, err
		}

		glog.Infof("%d NameUsages from %s", len(usages), srcType)

		if err := snowball.AddUsage(usages...); err != nil {
			return nil, err
		}

		glog.Infof("Current Aggregated Usages [%d]", snowball.Count())

	}

	return &snowball, nil
}

func fetchOccurrences(cxt context.Context, aggregation *aggregate.Aggregate) (*occurrence.Aggregation, error) {

	res := occurrence.Aggregation{}

	glog.Infof("Fetching Occurrences for %d NameUsages", aggregation.Count())

	// Fetch Inaturalist and GBIF occurrences.
	if err := aggregation.Each(cxt, occurrenceFetchFunc(&res, datasources.TypeINaturalist, datasources.TypeGBIF)); err != nil {
		return nil, err
	}

	filteredAggregation, err := aggregation.Filter(func(u nameusage.NameUsage) (bool, error) {
		i, err := u.Occurrences()
		if err != nil {
			return false, err
		}
		return i < 100, nil
	})
	if err != nil {
		return nil, err
	}

	sciNames, err := filteredAggregation.ScientificNames()
	if err != nil {
		return nil, err
	}

	glog.Info("Scientific Names from Occurrence Filtered Aggregation", sciNames)

	usages, err := sourcefetchers.FetchNameUsages(cxt, datasources.TypeMushroomObserver, sciNames, nil)
	if err != nil {
		return nil, err
	}

	glog.Infof("%d Usages from MushroomObserver", len(usages))

	if err := filteredAggregation.AddUsage(usages...); err != nil {
		return nil, err
	}

	if err := aggregation.Each(cxt, occurrenceFetchFunc(&res, datasources.TypeMushroomObserver)); err != nil {
		return nil, err
	}

	return &res, nil

}

func occurrenceFetchFunc(oAggr *occurrence.Aggregation, sourceTypes ...datasources.SourceType) aggregate.EachFunction {
	return func(ctx context.Context, usage nameusage.NameUsage) error {

		aggr, err := occurrence.FetchOccurrences(ctx, usage, false, sourceTypes...)
		if err != nil {
			return err
		}

		glog.Infof("Fetched %d Occurrences for NameUsage [%s]", aggr.Count(), usage.CanonicalName())

		if aggr.Count() >= minimumOccurrenceCount {
			return oAggr.Merge(aggr)
		}

		return nil
	}
}

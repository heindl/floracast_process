package main

import (
	"bitbucket.org/heindl/taxa/datasources/inaturalist"
	"context"
	"bitbucket.org/heindl/taxa/datasources/gbif"
	"bitbucket.org/heindl/taxa/nameusage/nameusage"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/datasources/natureserve"
	"bitbucket.org/heindl/taxa/datasources/mushroomobserver"
	"bitbucket.org/heindl/taxa/nameusage/aggregate"
	"bitbucket.org/heindl/taxa/occurrences"
	"flag"
	"strings"
	"strconv"
	"bitbucket.org/heindl/taxa/store"
)

const MinimumOccurrenceCount = 100


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

	aggr, err := InitialAggregation(cxt, intIDs...)
	if err != nil {
		panic(err)
	}

	ocAggr, err := OccurrenceFetch(cxt, aggr)
	if err != nil {
		panic(err)
	}

	fc, err := store.NewLiveFirestore()
	if err != nil {
		panic(err)
	}

	if err := ocAggr.Upsert(cxt, fc); err != nil {
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

func InitialAggregation(cxt context.Context, inaturalistTaxonIDs ...int) (*aggregate.Aggregate, error) {

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. TotalOccurrenceCount are still stored within the synonym.

	usages, err := inaturalist.FetchNameUsages(cxt, inaturalistTaxonIDs...)
	if err != nil {
		return nil, err
	}

	aggregation := aggregate.Aggregate{}
	if err := aggregation.AddUsage(usages...); err != nil {
		return nil, err
	}

	gbifUsages, err := gbif.FetchNamesUsages(cxt, aggregation.ScientificNames(), aggregation.TargetIDs(datasources.DataSourceTypeGBIF))
	if err != nil {
		return nil, err
	}

	if err := aggregation.AddUsage(gbifUsages...); err != nil {
		return nil, err
	}

	natureServeUsages, err := natureserve.FetchNameUsages(cxt, aggregation.ScientificNames(), aggregation.TargetIDs(datasources.DataSourceTypeNatureServe))
	if err != nil {
		return nil, err
	}

	if err := aggregation.AddUsage(natureServeUsages...); err != nil {
		return nil, err
	}

	return &aggregation, nil
}

func OccurrenceFetch(cxt context.Context, aggregation *aggregate.Aggregate) (*occurrences.OccurrenceAggregation, error) {

	res := occurrences.OccurrenceAggregation{}

	// Fetch Inaturalist and GBIF occurrences.
	if err := aggregation.Each(cxt, occurrenceFetcher(&res, datasources.DataSourceTypeGBIF, datasources.DataSourceTypeINaturalist)); err != nil {
		return nil, err
	}

	countFilteredAggregation, err := aggregation.Filter(func(u *nameusage.NameUsage) bool {
		return u.TotalOccurrenceCount() < 100
	})
	if err != nil {
		return nil, err
	}

	// Fetch MushroomObserver sources.
	if err := countFilteredAggregation.Each(cxt, fetchMushroomObserverSources); err != nil {
		return nil, err
	}

	// Fetch MushroomObserver Occurrences.
	if err := countFilteredAggregation.Each(cxt, occurrenceFetcher(&res, datasources.DataSourceTypeMushroomObserver)); err != nil {
		return nil, err
	}

	return &res, nil

}

func fetchMushroomObserverSources(ctx context.Context, usage *nameusage.NameUsage) error {
	names := append(usage.Synonyms().ScientificNames(), usage.CanonicalName().ScientificName())
	sources, err := mushroomobserver.MatchCanonicalNames(ctx, names...)
	if err != nil {
		return err
	}
	return usage.AddSources(sources...)
}

func occurrenceFetcher(oAggr *occurrences.OccurrenceAggregation, srcTypes ...datasources.SourceType) aggregate.EachFunction {
	return func(ctx context.Context, usage *nameusage.NameUsage) error {

		usageOccurrenceAggr := occurrences.OccurrenceAggregation{}

		for _, src := range usage.Sources(srcTypes...) {
			srcOccurrenceAggregation, err := occurrences.FetchOccurrences(ctx, src.SourceType(), src.TargetID(), src.LastFetchedAt())
			if err != nil {
				return err
			}
			if err := src.RegisterOccurrenceFetch(srcOccurrenceAggregation.Count()); err != nil {
				return err
			}
			if err := oAggr.Merge(srcOccurrenceAggregation); err != nil {
				return err
			}
		}

		if usageOccurrenceAggr.Count() >= MinimumOccurrenceCount {
			if err := oAggr.Merge(&usageOccurrenceAggr); err != nil {
				return err
			}
		}

		return nil

	}
}
package aggregate

import (
	"bitbucket.org/heindl/taxa/datasources/inaturalist"
	"context"
	"bitbucket.org/heindl/taxa/datasources/gbif"
	"bitbucket.org/heindl/taxa/nameusage"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/occurrences/occurrencefetcher"
	"bitbucket.org/heindl/taxa/datasources/natureserve"
	"bitbucket.org/heindl/taxa/datasources/mushroomobserver"
	"fmt"
)

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

func AggregateNameUsages(cxt context.Context, inaturalistTaxonIDs ...int) (*nameusage.AggregateNameUsages, error) {

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. TotalOccurrenceCount are still stored within the synonym.

	usages, err := inaturalist.FetchNameUsages(cxt, inaturalistTaxonIDs...)
	if err != nil {
		return nil, err
	}

	gbifUsages, err := gbif.FetchNamesUsages(cxt, usages.CanonicalNames().Strings(), usages.TargetIDs(datasources.DataSourceTypeGBIF))
	if err != nil {
		return nil, err
	}

	if err := usages.CombineWith(gbifUsages); err != nil {
		return nil, err
	}

	natureServeUsages, err := natureserve.FetchNameUsages(cxt, usages.CanonicalNames().Strings(), usages.TargetIDs(datasources.DataSourceTypeNatureServe))
	if err != nil {
		return nil, err
	}

	if err := usages.CombineWith(natureServeUsages); err != nil {
		return nil, err
	}

	// Fetch Inaturalist and GBIF occurrences.
	if err := usages.ForEach(cxt, newOccurrenceFetcher(datasources.DataSourceTypeINaturalist, datasources.DataSourceTypeGBIF)); err != nil {
		return nil, err
	}

	fmt.Println("USAGE COUNT AFTER FETCH", usages.Count())

	filteredUsages, err := usages.Filter(func(u *nameusage.CanonicalNameUsage) bool {
		fmt.Println("TOTAL OCCURRENCE COUNT IN FILTER", u.TotalOccurrenceCount())
		return u.TotalOccurrenceCount() < 100
	})

	fmt.Println("USAGE COUNT AFTER FILTER", filteredUsages.Count())

	// Fetch MushroomObserver sources.
	if err := filteredUsages.ForEach(cxt, fetchMushroomObserverSources); err != nil {
		return nil, err
	}

	// Fetch MushroomObserver Occurrences.
	if err := filteredUsages.ForEach(cxt, newOccurrenceFetcher(datasources.DataSourceTypeMushroomObserver)); err != nil {
		return nil, err
	}

	if err := filteredUsages.ForEach(cxt, func(ctx context.Context, usage *nameusage.CanonicalNameUsage) error {
		fmt.Println(usage.CanonicalName().ScientificName(), usage.SourceCount(), usage.TotalOccurrenceCount())
		fmt.Println("SourceCount", usage.SourceCount())
		return nil
	}); err != nil {
		return nil, err
	}



	return nil, nil

}

func fetchMushroomObserverSources(ctx context.Context, usage *nameusage.CanonicalNameUsage) error {
	sources, err := mushroomobserver.MatchCanonicalNames(ctx, usage.ScientificNames().Strings()...)
	if err != nil {
		return err
	}
	return usage.AddSources(sources...)
}

func newOccurrenceFetcher(sourceTypes ...datasources.DataSourceType) nameusage.ForEachModifierFunction {
	return func(ctx context.Context, usage *nameusage.CanonicalNameUsage) error {
		if usage == nil {
			return nil
		}
		for _, src := range usage.Sources(sourceTypes...) {
			if !src.TargetID().Valid(src.SourceType()) {
				fmt.Println(fmt.Sprintf("Warning: Attempting to fetch occurrences for an invalid source [%s, %s, %s]", src.SourceType(), src.TargetID(), src.CanonicalName().ScientificName()))
				continue
			}
			if !datasources.HasDataSourceType(sourceTypes, src.SourceType()) {
				continue
			}
			occurrenceAggregation, err := occurrencefetcher.FetchOccurrences(ctx, src.SourceType(), src.TargetID(), src.LastFetchedAt())
			if err != nil {
				return err
			}
			for _, o := range occurrenceAggregation.OccurrenceList() {
				if err := usage.AddOccurrence(o); err != nil {
					return err
				}
			}
		}

		return nil
	}
}
package aggregate

import (
	"bitbucket.org/heindl/taxa/inaturalist"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/gbif"
	"bitbucket.org/heindl/taxa/nature_serve"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
)

// This function is intended to be a process of building a store taxon representation based on several
// sorted that accumulate taxonomic information. These names and synonyms can then be used to fetch occurrences
// from various sources using only the list of names.

// Note that this function aggressively combines synonyms at this point. If any of the three sources consider it a synonym,
// we accept it in order error on the side of more occurrences.

// Note that this should also be the materialized view of the taxon used for displaying in the interface.

// Initially don't worry about existing values in the database. This function runs alone each time, connects to multiple
// datasources, and then deletes those in the database that match.


// The main goal here is occurrences. In order to that we have to know which names are for the same species,
// so that sources with different occurrences for different species can be combined.


// TODO: One day convert all of this to a graph database. Which could even incorporate DNA records.
// In this scenario, all occurrences would need photos, and the photos would be used instead of names to connect them
// to phenotypical dna.


// TODO: Include a step to verify the ids for occurrences haven't changed.
// Could actually be done in the occurrence fetch step: https://www.inaturalist.org/taxon_changes


func AggregateNameUsages(cxt context.Context, inaturalistTaxonIDs ...int) (name_usage.AggregateNameUsages, error) {

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. OccurrenceCount are still stored within the synonym.

	inaturalistUsages, err := inaturalist.FetchNameUsages(cxt, inaturalistTaxonIDs...)
	if err != nil {
		return nil, err
	}

	gbifUsages, err := gbif.FetchNamesUsages(cxt, inaturalistUsages.Names(), inaturalistUsages.TargetIDs(store.DataSourceTypeGBIF))
	if err != nil {
		return nil, err
	}

	aggregate := append(inaturalistUsages, gbifUsages...)

	natureServeUsages, err := nature_serve.FetchNameUsages(cxt, aggregate.Names(), aggregate.TargetIDs(store.DataSourceTypeNatureServe))
	if err != nil {
		return nil, err
	}

	aggregate = append(aggregate, natureServeUsages...)

	return aggregate.Condense()

}
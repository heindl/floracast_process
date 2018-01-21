package taxa

import (
	"bitbucket.org/heindl/taxa/inaturalist"
	"github.com/jessevdk/go-flags"
	"golang.org/x/net/context"
	"flag"
	"strings"
	"strconv"
	"bitbucket.org/heindl/taxa/store"
)

// This function is intended to be a process of building a store taxon representation based on several
// sorted that accumulate taxonomic information. These names and synonyms can then be used to fetch occurrences
// from various sources using only the list of names.

// Note that this should also be the materialized view of the taxon used for displaying in the interface.


// Initially don't worry about existing values in the database. This function runs alone each time, connects to multiple
// datasources, and then deletes those in the database that match.


// The main goal here is occurrences. In order to that we have to know which names are for the same species,
// so that sources with different occurrences for different species can be combined.


// TODO: One day convert all of this to a graph database. Which could even incorporate DNA records.
// In this scenario, all occurrences would need photos, and the photos would be used instead of names to connect them
// to phenotypical dna.

func main() {

	iTaxaIDStr := flag.String("inaturalist_taxa_ids", "", "parent taxa for query, string separated")
	flag.Parse()

	iTaxaIDs := []inaturalist.TaxonID{}
	for _, id := range strings.Split(*iTaxaIDStr, ",") {
		inaturalistTaxaID, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		iTaxaIDs = append(iTaxaIDs, inaturalist.TaxonID(inaturalistTaxaID))
	}

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. Occurrences are still stored within the synonym.

	iTaxa, err := inaturalist.FetchTaxaAndChildren(context.Background(), iTaxaIDs...)
	if err != nil {
		panic(err)
	}

	inaturalistLinks := map[Name][]inaturalist.INaturalistTaxonScheme{}
	nameGroup := map[Name][]NameSource{}

	for _, iTaxon := range iTaxa {
		name := NewName(iTaxon.Name)
		sourceID := store.DataSourceIDINaturalist
		targetID := store.DataSourceTargetID(strconv.Itoa(int(iTaxon.ID)))

		if _, ok := inaturalistLinks[name]; !ok {
			inaturalistLinks[name] = []inaturalist.INaturalistTaxonScheme{}
		}
		inaturalistLinks[name] = append(inaturalistLinks[name], iTaxon.TaxonSchemes...)

		if _, ok := nameGroup[name]; !ok {
			nameGroup[name] = []NameSource{}
		}

		// TODO: Include process to add synonyms here, though right now we haven't found any.
		nameGroup[name] = append(nameGroup[name], NameSource{
			SourceID: sourceID,
			TargetID: targetID,
			IsSynonym: false,
		})
	}


	// Include a step to verify the ids for occurrences haven't changed.
	// Could actually be done in the occurrence fetch step: https://www.inaturalist.org/taxon_changes

}

type Name string

func NewName(s string) Name {
	return Name(s).Format()
}

func (n Name) Format() Name {
	return Name(strings.ToLower(string(n)))
}

type NameSource struct {
	SourceID store.DataSourceID
	TargetID store.DataSourceTargetID
	IsSynonym bool
	Synonyms []Name
}

type Taxon struct {
	Names map[Name][]NameSource
}

type orchestrator struct {
	Taxa []*Taxon
	INaturalistTaxa []*inaturalist.Taxon
	//GBIFTaxa []
	NatureServe []
}
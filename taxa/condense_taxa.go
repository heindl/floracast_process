package main

import (
	"bitbucket.org/heindl/taxa/inaturalist"
	"golang.org/x/net/context"
	"flag"
	"strings"
	"strconv"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/gbif"
	"bitbucket.org/heindl/taxa/nature_serve"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/utils"
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

func main() {
	iTaxaIDStr := flag.String("inaturalist_taxa_ids", "", "parent taxa for query, string separated")
	flag.Parse()

	_, err := fetchTaxaSources(context.Background(), inaturalist.ParseStringIDs(strings.Split(*iTaxaIDStr, ",")...)...)
	if err != nil {
		panic(err)
	}



}


func fetchTaxaSources(cxt context.Context, taxonIDs ...inaturalist.TaxonID) (CanonicalNameSources, error) {

	// Start with GBIF because the hiearchy is simple. The occurrence sources for the gbif will be searched externally.
	// Note also that Inaturalist appears to try to avoid synonyms: https://www.inaturalist.org/taxon_changes
	// Which means that try to combine them and ignore synonyms, though they appear to still show known synonyms like Morchella Conica.
	// We need synonyms because other archives do not appear to divide them. Occurrences are still stored within the synonym.

	inaturalistTaxa, err := inaturalist.FetchTaxaAndChildren(cxt, taxonIDs...)
	if err != nil {
		return nil, err
	}

	sources, gbifKeyMap, err := parseINaturalistTaxa(cxt, inaturalistTaxa...)
	if err != nil {
		return nil, err
	}

	//fmt.Println(utils.JsonOrSpew(sources))
	//fmt.Println("m := map[int][]string{")
	//for k, names := range gbifKeyMap {
	//	fmt.Println(fmt.Sprintf(`%d: []string{`, k))
	//	for _, n := range names {
	//		fmt.Println(fmt.Sprintf(`"%s",`, n))
	//	}
	//	fmt.Println(`},`)
	//}
	//fmt.Println("}")
	//
	//return nil, nil

	natureServeTaxa, err := nature_serve.FetchTaxaFromSearch(cxt, sources.Names()...)
	if err != nil {
		return nil, err
	}
	sources = append(sources, parseNatureServeTaxa(natureServeTaxa...)...)

	gbifUsages, err := gbif.FetchNamesUsages(cxt, sources.Names(), gbifKeyMap)
	if err != nil {
		return nil, err
	}
	sources = append(sources, parseGBIFTaxa(gbifUsages)...)


	return sources, nil

}

func parseINaturalistTaxa(cxt context.Context, taxa ...*inaturalist.INaturalistTaxon) (CanonicalNameSources, map[int][]string, error) {

	sources := CanonicalNameSources{}

	gbifKeyMap := map[int][]string{}

	for _, inaturalistTaxon := range taxa {

		id := store.DataSourceTargetID(strconv.Itoa(int(inaturalistTaxon.ID)))

		// TODO: Currently no way to handle synonym, but also no synonyms coming from inaturalist.
		sources = append(sources, CanonicalNameSource{
			CanonicalName:     strings.ToLower(inaturalistTaxon.Name),
			SourceOccurrences: SourceOccurrenceCount{store.DataSourceIDINaturalist: map[store.DataSourceTargetID]int{id: inaturalistTaxon.ObservationsCount}},
			Ranks:             []string{strings.ToLower(inaturalistTaxon.Rank)},
		})

		for _, scheme := range inaturalistTaxon.TaxonSchemes {
			if scheme.DataSourceID == store.DataSourceIDGBIF {
				gbifKey, err := strconv.Atoi(string(scheme.TargetID))
				if err != nil {
					return nil, nil, errors.Wrapf(err, "could not parse target id [%s]", scheme.TargetID)
				}
				if _, ok := gbifKeyMap[gbifKey]; !ok {
					gbifKeyMap[gbifKey] = []string{}
				}
				gbifKeyMap[gbifKey] = append(gbifKeyMap[gbifKey], inaturalistTaxon.Name)
			}
			if scheme.DataSourceID == store.DataSourceNatureServe {
				natureServeTxn, err := nature_serve.FetchTaxonWithUID(cxt, string(scheme.TargetID), inaturalistTaxon.Name)
				if err != nil {
					return nil, nil, err
				}
				sources = append(sources, parseNatureServeTaxa(natureServeTxn)...)
			}
		}
	}

	return sources, gbifKeyMap, nil
}

func parseGBIFTaxa(usages gbif.CanonicalNameUsages) CanonicalNameSources {
	srcs := CanonicalNameSources{}

	for _, txn := range usages {
		occurrences := map[store.DataSourceTargetID]int{}
		for id, count := range txn.OccurrenceCount {
			occurrences[store.DataSourceTargetID(strconv.Itoa(int(id)))] = count
		}

		src := CanonicalNameSource{
			CanonicalName:     strings.ToLower(txn.CanonicalName),
			SourceOccurrences: SourceOccurrenceCount{store.DataSourceIDGBIF: occurrences},
			Ranks:             txn.Ranks,
			Synonyms: txn.Synonyms,
		}

		srcs = append(srcs, src)

	}

	return srcs
}

func parseNatureServeTaxa(taxa ...*nature_serve.Taxon) CanonicalNameSources {
	srcs := CanonicalNameSources{}

	for _, txn := range taxa {

		src := CanonicalNameSource{
			CanonicalName: strings.ToLower(txn.ScientificName.Name),
			SourceOccurrences: SourceOccurrenceCount{store.DataSourceNatureServe: map[store.DataSourceTargetID]int{
				store.DataSourceTargetID(txn.ID): 0,
			}},
			Ranks: []string{"species"},
		}

		for _, nsSynonym := range txn.Synonyms {
			src.Synonyms = utils.AddStringToSet(src.Synonyms, strings.ToLower(nsSynonym.Name))
			//src.SourceOccurrenceCount[store.DataSourceNatureServe] = append(
			//	src.SourceOccurrenceCount[store.DataSourceNatureServe],
			//	store.DataSourceTargetID(nsSynonym.ConceptReferenceCode),
			//)
		}

		srcs = append(srcs, src)

	}

	return srcs
}
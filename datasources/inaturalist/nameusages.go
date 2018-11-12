package inaturalist

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/nameusage/canonicalname"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"context"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

// FetchNameUsages fetches usages from known INaturalist IDs.
func FetchNameUsages(cxt context.Context, _ []string, ids datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	ints, err := ids.Integers()
	if err != nil {
		return nil, err
	}
	taxonIDs := taxonIDsFromIntegers(ints...)

	taxa, err := newTaxaFetcher(cxt, true, true).FetchTaxa(taxonIDs...)
	if err != nil {
		return nil, err
	}

	res := []nameusage.NameUsage{}

	for _, inaturalistTaxon := range taxa {
		usage, err := parseNameUsage(inaturalistTaxon)
		if err != nil {
			return nil, err
		}
		res = append(res, usage)
	}

	return res, nil
}

func parseNameUsage(txn *taxon) (nameusage.NameUsage, error) {
	if len(txn.CurrentSynonymousTaxonIds) > 0 {
		// So far this has never been the case, but if it is, we need to process those.
		return nil, errors.Newf("INaturalist taxon [%d] has synonymous taxon IDs [%v]", txn.ID, txn.CurrentSynonymousTaxonIds)
	}

	name, err := canonicalname.NewCanonicalName(txn.Name, strings.ToLower(txn.Rank))
	if err != nil {
		return nil, err
	}

	src, err := nameusage.NewSource(datasources.TypeINaturalist, txn.ID.TargetID(), name)
	if err != nil {
		return nil, err
	}

	if txn.PreferredCommonName != "" {
		if err = src.AddCommonNames(txn.PreferredCommonName); err != nil {
			return nil, err
		}
	}

	usage, err := nameusage.NewNameUsage(src)
	if err != nil {
		return nil, err
	}

	for _, scheme := range txn.TaxonSchemes {

		src, err := nameusage.NewSource(scheme.SourceType, scheme.TargetID, name)
		if err != nil {
			return nil, err
		}

		if err := usage.AddSources(src); err != nil {
			return nil, err
		}
	}
	return usage, nil
}

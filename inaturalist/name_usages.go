package inaturalist

import (
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"bitbucket.org/heindl/taxa/store"
	"context"
	"strings"
	"github.com/saleswise/errors/errors"
)

func FetchNameUsages(cxt context.Context, ids ...int) (name_usage.CanonicalNameUsages, error) {


	taxa, err := FetchTaxaAndChildren(cxt, TaxonIDsFromIntegers(ids...)...)
	if err != nil {
		return nil, err
	}

	res := name_usage.CanonicalNameUsages{}

	for _, inaturalistTaxon := range taxa {
		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			return nil, errors.Newf("Unexpected synonyms [%+v] from INaturalist taxon [%d]", inaturalistTaxon.CurrentSynonymousTaxonIds, inaturalistTaxon.ID)
		}
		usage := name_usage.CanonicalNameUsage{
			CanonicalName:     strings.ToLower(inaturalistTaxon.Name),
			SourceTargetOccurrenceCount: name_usage.SourceTargetOccurrenceCount{},
			Ranks:             []string{strings.ToLower(inaturalistTaxon.Rank)},
		}

		usage.SourceTargetOccurrenceCount.Set(store.DataSourceTypeINaturalist, inaturalistTaxon.ID.TargetID(), inaturalistTaxon.ObservationsCount)
		for _, scheme := range inaturalistTaxon.TaxonSchemes {
			usage.SourceTargetOccurrenceCount.Set(scheme.DataSourceID, scheme.TargetID, 0)
		}
		res = append(res, usage)
	}

	return res.Condense()
}

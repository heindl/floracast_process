package inaturalist

import (
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"bitbucket.org/heindl/taxa/store"
	"context"
	"strings"
	"github.com/saleswise/errors/errors"
)

func FetchNameUsages(cxt context.Context, ids ...int) (*name_usage.AggregateNameUsages, error) {

	taxa, err := FetchTaxaAndChildren(cxt, TaxonIDsFromIntegers(ids...)...)
	if err != nil {
		return nil, err
	}

	res := name_usage.AggregateNameUsages{}

	for _, inaturalistTaxon := range taxa {
		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			return nil, errors.Newf("Unexpected synonyms [%+v] from INaturalist taxon [%d]", inaturalistTaxon.CurrentSynonymousTaxonIds, inaturalistTaxon.ID)
		}

		name, err := name_usage.NewCanonicalName(inaturalistTaxon.Name, strings.ToLower(inaturalistTaxon.Rank))
		if err != nil {
			return nil, err
		}

		src, err := name_usage.NewNameUsageSource(store.DataSourceTypeINaturalist, inaturalistTaxon.ID.TargetID(), name, true)
		if err != nil {
			return nil, err
		}

		usage, err := name_usage.NewCanonicalNameUsage(src)
		if err != nil {
			return nil, err
		}

		for _, scheme := range inaturalistTaxon.TaxonSchemes {

			src, err := name_usage.NewNameUsageSource(scheme.DataSourceType, scheme.TargetID, name, true)
			if err != nil {
				return nil, err
			}

			if err := usage.AddSource(src); err != nil {
				return nil, err
			}
		}

		if err := res.Add(usage); err != nil {
			return nil, err
		}

	}

	return &res, nil
}

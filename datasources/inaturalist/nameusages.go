package inaturalist

import (
	"bitbucket.org/heindl/taxa/nameusage"
	"context"
	"strings"
	"github.com/saleswise/errors/errors"
	"bitbucket.org/heindl/taxa/datasources"
)

func FetchNameUsages(cxt context.Context, ids ...int) (*nameusage.AggregateNameUsages, error) {

	taxa, err := FetchTaxaAndChildren(cxt, TaxonIDsFromIntegers(ids...)...)
	if err != nil {
		return nil, err
	}

	res := nameusage.AggregateNameUsages{}

	for _, inaturalistTaxon := range taxa {
		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			return nil, errors.Newf("Unexpected synonyms [%+v] from INaturalist taxon [%d]", inaturalistTaxon.CurrentSynonymousTaxonIds, inaturalistTaxon.ID)
		}

		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			// So far this has never been the case, but if it is, we need to process those.
			return nil, errors.Newf("Taxon [%d] has synonymous Taxon IDs [%v]", inaturalistTaxon.ID, inaturalistTaxon.CurrentSynonymousTaxonIds)
		}

		name, err := nameusage.NewCanonicalName(inaturalistTaxon.Name, strings.ToLower(inaturalistTaxon.Rank))
		if err != nil {
			return nil, err
		}

		src, err := nameusage.NewNameUsageSource(datasources.DataSourceTypeINaturalist, inaturalistTaxon.ID.TargetID(), name)
		if err != nil {
			return nil, err
		}

		if inaturalistTaxon.PreferredCommonName != "" {
			if err := src.AddCommonNames(inaturalistTaxon.PreferredCommonName); err != nil {
				return nil, err
			}
		}

		usage, err := nameusage.NewCanonicalNameUsage(src)
		if err != nil {
			return nil, err
		}

		for _, scheme := range inaturalistTaxon.TaxonSchemes {

			src, err := nameusage.NewNameUsageSource(scheme.DataSourceType, scheme.TargetID, name)
			if err != nil {
				return nil, err
			}

			if err := usage.AddSources(src); err != nil {
				return nil, err
			}
		}

		if err := res.AddUsages(usage); err != nil {
			return nil, err
		}

	}

	return &res, nil
}

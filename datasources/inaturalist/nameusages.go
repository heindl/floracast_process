package inaturalist

import (
	"bitbucket.org/heindl/taxa/nameusage/nameusage"
	"context"
	"strings"
	"github.com/saleswise/errors/errors"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/nameusage/canonicalname"
	"bitbucket.org/heindl/taxa/nameusage/nameusagesource"
)

func FetchNameUsages(cxt context.Context, ids ...int) ([]*nameusage.NameUsage, error) {

	taxa, err := NewTaxaFetcher(cxt, true, true).FetchTaxa(TaxonIDsFromIntegers(ids...)...)
	if err != nil {
		return nil, err
	}

	res := []*nameusage.NameUsage{}

	for _, inaturalistTaxon := range taxa {
		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			return nil, errors.Newf("Unexpected synonyms [%+v] from INaturalist taxon [%d]", inaturalistTaxon.CurrentSynonymousTaxonIds, inaturalistTaxon.ID)
		}

		if len(inaturalistTaxon.CurrentSynonymousTaxonIds) > 0 {
			// So far this has never been the case, but if it is, we need to process those.
			return nil, errors.Newf("Taxon [%d] has synonymous Taxon IDs [%v]", inaturalistTaxon.ID, inaturalistTaxon.CurrentSynonymousTaxonIds)
		}

		name, err := canonicalname.NewCanonicalName(inaturalistTaxon.Name, strings.ToLower(inaturalistTaxon.Rank))
		if err != nil {
			return nil, err
		}

		src, err := nameusagesource.NewSource(datasources.DataSourceTypeINaturalist, inaturalistTaxon.ID.TargetID(), name)
		if err != nil {
			return nil, err
		}

		if inaturalistTaxon.PreferredCommonName != "" {
			if err := src.AddCommonNames(inaturalistTaxon.PreferredCommonName); err != nil {
				return nil, err
			}
		}

		usage, err := nameusage.NewNameUsage(src)
		if err != nil {
			return nil, err
		}

		for _, scheme := range inaturalistTaxon.TaxonSchemes {

			src, err := nameusagesource.NewSource(scheme.SourceType, scheme.TargetID, name)
			if err != nil {
				return nil, err
			}

			if err := usage.AddSources(src); err != nil {
				return nil, err
			}
		}

		res = append(res, usage)


	}

	return res, nil
}

package natureserve

import (
	"context"
	"bitbucket.org/heindl/taxa/nameusage"
	"bitbucket.org/heindl/taxa/datasources"
)

func FetchNameUsages(cxt context.Context, names []string, targetIDs datasources.DataSourceTargetIDs) (*nameusage.AggregateNameUsages, error) {

	nameTaxa, err := FetchTaxaFromSearch(cxt, names...)
	if err != nil {
		return nil, err
	}

	uidTaxa, err := FetchTaxaWithUID(cxt, targetIDs.Strings()...)
	if err != nil {
		return nil, err
	}

	taxa := append(nameTaxa, uidTaxa...)

	res := nameusage.AggregateNameUsages{}

	for _, txn := range taxa {


		canonicalName, err := nameusage.NewCanonicalName(txn.ScientificName.Name, "species")
		if err != nil {
			return nil, err
		}

		usageSource, err := nameusage.NewNameUsageSource(datasources.DataSourceTypeNatureServe, datasources.DataSourceTargetID(txn.ID), canonicalName)
		if err != nil {
			return nil, err
		}

		for _, commonName := range txn.CommonNames {
			if err := usageSource.AddCommonNames(commonName.Name); err != nil {
				return nil, err
			}
		}

		for _, synonym := range txn.Synonyms {
			synonymCanonicalName, err := nameusage.NewCanonicalName(synonym.Name, "species")
			if err != nil {
				return nil, err
			}
			if err := usageSource.AddSynonym(synonymCanonicalName); err != nil {
				return nil, err
			}
		}

		usage, err := nameusage.NewCanonicalNameUsage(usageSource)
		if err != nil {
			return nil, err
		}


		if err := res.AddUsages(usage); err != nil {
			return nil, err
		}
	}

	return &res, nil

}
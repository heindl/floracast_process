package natureserve

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"context"
)

func FetchNameUsages(cxt context.Context, names []string, targetIDs datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	nameTaxa, err := FetchTaxaFromSearch(cxt, names...)
	if err != nil {
		return nil, err
	}

	uidTaxa, err := FetchTaxaWithUID(cxt, targetIDs.Strings()...)
	if err != nil {
		return nil, err
	}

	taxa := append(nameTaxa, uidTaxa...)

	res := []nameusage.NameUsage{}

	for _, txn := range taxa {

		canonicalName, err := canonicalname.NewCanonicalName(txn.ScientificName.Name, "species")
		if err != nil {
			return nil, err
		}

		usageSource, err := nameusage.NewSource(datasources.TypeNatureServe, datasources.TargetID(txn.ID), canonicalName)
		if err != nil {
			return nil, err
		}

		for _, commonName := range txn.CommonNames {
			if err := usageSource.AddCommonNames(commonName.Name); err != nil {
				return nil, err
			}
		}

		for _, synonym := range txn.Synonyms {
			synonymCanonicalName, err := canonicalname.NewCanonicalName(synonym.Name, "species")
			if err != nil {
				return nil, err
			}
			if err := usageSource.AddSynonym(synonymCanonicalName); err != nil {
				return nil, err
			}
		}

		usage, err := nameusage.NewNameUsage(usageSource)
		if err != nil {
			return nil, err
		}

		res = append(res, usage)

	}

	return res, nil

}

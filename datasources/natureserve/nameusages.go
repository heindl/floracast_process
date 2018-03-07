package natureserve

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"context"
)

// FetchNameUsages implements the NameUsage fetch interface.
func FetchNameUsages(cxt context.Context, names []string, targetIDs datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	nameTaxa, err := fetchTaxaFromSearch(cxt, names...)
	if err != nil {
		return nil, err
	}

	uidTaxa, err := fetchTaxaWithUID(cxt, targetIDs.Strings()...)
	if err != nil {
		return nil, err
	}

	taxa := append(nameTaxa, uidTaxa...)

	res := []nameusage.NameUsage{}

	for _, txn := range taxa {
		usage, err := txn.asNameUsage()
		if err != nil {
			return nil, err
		}
		res = append(res, usage)
	}

	return res, nil

}

func (Ω *taxon) asNameUsage() (nameusage.NameUsage, error) {
	canonicalName, err := canonicalname.NewCanonicalName(Ω.ScientificName.Name, "species")
	if err != nil {
		return nil, err
	}

	usageSource, err := nameusage.NewSource(datasources.TypeNatureServe, datasources.TargetID(Ω.ID), canonicalName)
	if err != nil {
		return nil, err
	}

	for _, commonName := range Ω.CommonNames {
		if err := usageSource.AddCommonNames(commonName.Name); err != nil {
			return nil, err
		}
	}

	for _, synonym := range Ω.Synonyms {
		synonymCanonicalName, err := canonicalname.NewCanonicalName(synonym.Name, "species")
		if err != nil {
			return nil, err
		}
		if err := usageSource.AddSynonym(synonymCanonicalName); err != nil {
			return nil, err
		}
	}
	return nameusage.NewNameUsage(usageSource)
}

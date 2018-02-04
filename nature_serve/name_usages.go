package nature_serve

import (
	"context"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"strings"
	"bitbucket.org/heindl/taxa/utils"
)

func FetchNameUsages(cxt context.Context, names []string, targetIDs store.DataSourceTargetIDs) (name_usage.AggregateNameUsages, error) {

	nameTaxa, err := FetchTaxaFromSearch(cxt, names...)
	if err != nil {
		return nil, err
	}

	uidTaxa, err := FetchTaxaWithUID(cxt, targetIDs.Strings()...)
	if err != nil {
		return nil, err
	}

	taxa := append(nameTaxa, uidTaxa...)

	res := name_usage.AggregateNameUsages{}

	for _, txn := range taxa {

		src := name_usage.CanonicalNameUsage{
			CanonicalName: strings.ToLower(txn.ScientificName.Name),
			SourceTargetOccurrenceCount: name_usage.sourceTargetOccurrenceCount{},
			Ranks: []string{"species"},
		}
		src.SourceTargetOccurrenceCount.Set(store.DataSourceTypeNatureServe, store.DataSourceTargetID(txn.ID), 0)

		// TODO: Why wouldn't these synonyms have an id attached to them. Should investigate when time.
		for _, nsSynonym := range txn.Synonyms {
			src.Synonyms = utils.AddStringToSet(src.Synonyms, strings.ToLower(nsSynonym.Name))
		}

		res = append(res, src)
	}

	return res.Condense()

}
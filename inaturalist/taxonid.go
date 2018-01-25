package inaturalist

import (
	"strconv"
	"bitbucket.org/heindl/taxa/store"
)

type TaxonID int64

func (Ω TaxonID) Valid() bool {
	return Ω != 0
}

func TaxonIDFromTargetID(id store.DataSourceTargetID) TaxonID {
	i, err := strconv.Atoi(string(id))
	if err != nil {
		return TaxonID(0)
	}
	return TaxonID(i)
}

func (Ω TaxonID) TargetID() store.DataSourceTargetID {
	return store.DataSourceTargetID(strconv.Itoa(int(Ω)))
}

func TaxonIDsFromIntegers(ids ...int) (res []TaxonID) {
	for _, id := range ids {
		res = append(res, TaxonID(id))
	}
	return
}

type TaxonIDs []TaxonID

func (Ω TaxonIDs) IndexOf(id TaxonID) int {
	for i := range Ω {
		if Ω[i] == id {
			return i
		}
	}
	return -1
}

func (Ω TaxonIDs) AddToSet(id TaxonID) TaxonIDs {
	if Ω.IndexOf(id) == -1 {
		return append(Ω, id)
	}
	return Ω
}
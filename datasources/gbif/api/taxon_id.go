package api

import (
	"strconv"
	"bitbucket.org/heindl/processors/datasources"
)

type TaxonID int

func (Ω TaxonID) Valid() bool {
	return Ω != 0
}

func (Ω TaxonID) TargetID() datasources.TargetID {
	return datasources.TargetID(strconv.Itoa(int(Ω)))
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
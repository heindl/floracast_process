package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"strconv"
)

type taxonID int64

func (Ω taxonID) Valid() bool {
	return Ω != 0
}

func taxonIDFromTargetID(id datasources.TargetID) taxonID {
	i, err := strconv.Atoi(string(id))
	if err != nil {
		return taxonID(0)
	}
	return taxonID(i)
}

func (Ω taxonID) TargetID() datasources.TargetID {
	return datasources.TargetID(strconv.Itoa(int(Ω)))
}

func taxonIDsFromIntegers(ids ...int) (res []taxonID) {
	for _, id := range ids {
		res = append(res, taxonID(id))
	}
	return
}

type taxonIDs []taxonID

func (Ω taxonIDs) IndexOf(id taxonID) int {
	for i := range Ω {
		if Ω[i] == id {
			return i
		}
	}
	return -1
}

func (Ω taxonIDs) AddToSet(id taxonID) taxonIDs {
	if Ω.IndexOf(id) == -1 {
		return append(Ω, id)
	}
	return Ω
}

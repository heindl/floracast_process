package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"strconv"
)

type taxonID int

func (Ω taxonID) Valid() bool {
	return Ω != 0
}

func (Ω taxonID) TargetID() datasources.TargetID {
	return datasources.TargetID(strconv.Itoa(int(Ω)))
}

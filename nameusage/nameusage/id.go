package nameusage

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/elgs/gostrgen"
	"math"
	"strings"
)

// https://github.com/cheekybits/genny -> Look into generation

type NameUsageID string
type NameUsageIDs []NameUsageID

func (Ω NameUsageID) Valid() bool {
	return len(Ω.String()) == nameUsageIDLength
}

func (Ω NameUsageID) String() string {
	return strings.TrimSpace(string(Ω))
}

func (Ω NameUsageIDs) AddToSet(æ NameUsageID) NameUsageIDs {
	if Ω.Contains(æ) {
		return Ω
	}
	return append(Ω, æ)
}

func (Ω NameUsageIDs) Contains(æ NameUsageID) bool {
	for _, id := range Ω {
		if id.String() == æ.String() {
			return true
		}
	}
	return false
}

func (Ω NameUsageIDs) Batch(maxBatchSize float64) []NameUsageIDs {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []NameUsageIDs{}
	for i := 0.0; i <= batchCount-1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)

		if end > len(Ω) {
			end = len(Ω)
		}
		o := Ω[start:end]
		res = append(res, o)
	}

	return res
}

func NameUsageIDsFromStrings(æ []string) (NameUsageIDs, error) {
	res := NameUsageIDs{}
	for _, strID := range æ {
		id := NameUsageID(strID)
		if !id.Valid() {
			return nil, errors.Newf("Invalid NameUsageID [%s]", id)
		}
		res = append(res, id)
	}
	return res, nil
}

const nameUsageIDLength = 25

func newNameUsageID() (NameUsageID, error) {
	rand, err := gostrgen.RandGen(nameUsageIDLength, gostrgen.Lower|gostrgen.Digit|gostrgen.Upper, "", "")
	if err != nil {
		return "", errors.Wrap(err, "Could not generate name usage id")
	}
	id := NameUsageID(rand)
	if !id.Valid() {
		return NameUsageID(""), errors.Newf("Creating invalid NameUsageID [%s]", id)
	}
	return id, nil
}

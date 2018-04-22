package nameusage

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/elgs/gostrgen"
	"math"
	"strings"
)

// https://github.com/cheekybits/genny -> Look into generation

// ID is an automatically generated ID used in FireStore
type ID string

// IDs is a list of IDs
type IDs []ID

// Valid ...
func (Ω ID) Valid() bool {
	return len(Ω.String()) == nameUsageIDLength
}

// String ...
func (Ω ID) String() string {
	return strings.TrimSpace(string(Ω))
}

// AddToSet ...
func (Ω IDs) AddToSet(æ ID) IDs {
	if Ω.Contains(æ) {
		return Ω
	}
	return append(Ω, æ)
}

// Contains ...
func (Ω IDs) Contains(æ fmt.Stringer) bool {
	for _, id := range Ω {
		if id.String() == æ.String() {
			return true
		}
	}
	return false
}

// Batch returns an array of arrays
func (Ω IDs) Batch(maxBatchSize float64) []IDs {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []IDs{}
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

// IDsFromStrings validates and casts strings as IDs.
func IDsFromStrings(æ []string) (IDs, error) {
	res := IDs{}
	for _, strID := range æ {
		id := ID(strID)
		if !id.Valid() {
			return nil, errors.Newf("Invalid ID [%s]", id)
		}
		res = res.AddToSet(id)
	}
	return res, nil
}

const nameUsageIDLength = 7

func NewNameUsageID() (ID, error) {
	rand, err := gostrgen.RandGen(nameUsageIDLength, gostrgen.Lower|gostrgen.Digit|gostrgen.Upper, "", "")
	if err != nil {
		return "", errors.Wrap(err, "Could not generate name usage id")
	}
	id := ID(rand)
	if !id.Valid() {
		return ID(""), errors.Newf("Creating invalid ID [%s]", id)
	}
	return id, nil
}

package datasources

import (
	"strconv"
	"github.com/dropbox/godropbox/errors"
)

type TargetID string

func (Ω TargetID) Valid(sourceType SourceType) bool {

	if  string(Ω) == "" || string(Ω) == "0" {
		return false
	}

	intTypes := []SourceType{TypeGBIF, TypeINaturalist, TypeMushroomObserver}

	_, intParseErr := strconv.Atoi(string(Ω))

	if intParseErr != nil && HasDataSourceType(intTypes, sourceType) {
		return false
	}

	if intParseErr == nil && HasDataSourceType([]SourceType{TypeNatureServe}, sourceType) {
		return false
	}

	return true
}

func (Ω TargetID) ToInt() (int, error) {
	i, err := strconv.Atoi(string(Ω))
	if err != nil {
		return 0, errors.Wrapf(err, "Could not cast TargetID [%s] as int", Ω)
	}
	return i, nil
}

func NewDataSourceTargetIDFromInt(i int) (TargetID, error) {
	if i == 0 {
		return TargetID(""), errors.New("Invalid TargetID: Received zero.")
	}
	return TargetID(strconv.Itoa(i)), nil
}

func NewDataSourceTargetIDFromInts(ints ...int) (TargetIDs, error) {
	res := TargetIDs{}
	for _, i := range ints {
		id, err := NewDataSourceTargetIDFromInt(i)
		if err != nil {
			return nil, err
		}
		res = res.AddToSet(id)
	}
	return res, nil
}

type TargetIDs []TargetID

func (Ω TargetIDs) Integers() ([]int, error) {
	res := []int{}
	for _, id := range Ω {
		i, err := id.ToInt()
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (Ω TargetIDs) Strings() (res []string) {
	for _, id := range Ω {
		res = append(res, string(id))
	}
	return
}

func (Ω TargetIDs) AddToSet(ids ...TargetID) TargetIDs {
	for _, id := range ids {
		if Ω.Contains(id) {
			continue
		}
		Ω = append(Ω, id)
	}
	return Ω
}

func (Ω TargetIDs) Contains(id TargetID) bool {
	for i := range Ω {
		if Ω[i] == id {
			return true
		}
	}
	return false
}

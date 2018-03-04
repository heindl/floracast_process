package datasources

import (
	"github.com/dropbox/godropbox/errors"
	"strconv"
	"strings"
)

type TargetID string

type TargetIDProvider func() (TargetID, error)

func NewTargetID(target string, sourceType SourceType) (TargetID, error) {
	targetID := TargetID(target)
	if !targetID.Valid(sourceType) {
		return TargetID(""), errors.Newf("Invalid TargetID [%s] with SourceType [%s]", target, sourceType)
	}
	return targetID, nil
}

func (Ω TargetID) Valid(sourceType SourceType) bool {

	if !sourceType.Valid() {
		return false
	}

	intSourceTypes := []SourceType{TypeGBIF, TypeINaturalist, TypeMushroomObserver}
	strSourceTypes := []SourceType{TypeNatureServe}

	s := strings.TrimSpace(string(Ω))

	if s == "" || s == "0" {
		return false
	}

	// Allow unchecked ones to fall through
	if !HasDataSourceType(append(strSourceTypes, intSourceTypes...), sourceType) {
		return true
	}

	_, intParseErr := strconv.Atoi(s)

	if intParseErr != nil && HasDataSourceType(intSourceTypes, sourceType) {
		return false
	}

	if intParseErr == nil && HasDataSourceType(strSourceTypes, sourceType) {
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

func NewDataSourceTargetIDFromInt(sourceType SourceType, i int) (TargetID, error) {
	return NewTargetID(strconv.Itoa(i), sourceType)
}

func NewDataSourceTargetIDFromInts(sourceType SourceType, ints ...int) (TargetIDs, error) {
	res := TargetIDs{}
	for _, i := range ints {
		id, err := NewDataSourceTargetIDFromInt(sourceType, i)
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

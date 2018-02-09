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

	intTypes := []SourceType{DataSourceTypeGBIF, DataSourceTypeINaturalist, DataSourceTypeMushroomObserver}

	_, intParseErr := strconv.Atoi(string(Ω))

	if intParseErr != nil && HasDataSourceType(intTypes, sourceType) {
		return false
	}

	if intParseErr == nil && HasDataSourceType([]SourceType{DataSourceTypeNatureServe}, sourceType) {
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

type DataSourceTargetIDs []TargetID

func (Ω DataSourceTargetIDs) Strings() (res []string) {
	for _, id := range Ω {
		res = append(res, string(id))
	}
	return
}

func (Ω DataSourceTargetIDs) AddToSet(ids ...TargetID) DataSourceTargetIDs {
	for _, id := range ids {
		if Ω.Contains(id) {
			continue
		}
		Ω = append(Ω, id)
	}
	return Ω
}

func (Ω DataSourceTargetIDs) Contains(id TargetID) bool {
	for i := range Ω {
		if Ω[i] == id {
			return true
		}
	}
	return false
}

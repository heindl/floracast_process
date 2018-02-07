package datasources

import (
	"strconv"
	"github.com/dropbox/godropbox/errors"
)

type DataSourceTargetID string

func (Ω DataSourceTargetID) Valid(sourceType DataSourceType) bool {

	if  string(Ω) == "" || string(Ω) == "0" {
		return false
	}

	intTypes := []DataSourceType{DataSourceTypeGBIF, DataSourceTypeINaturalist, DataSourceTypeMushroomObserver}

	_, intParseErr := strconv.Atoi(string(Ω))

	if intParseErr != nil && HasDataSourceType(intTypes, sourceType) {
		return false
	}

	if intParseErr == nil && HasDataSourceType([]DataSourceType{DataSourceTypeMushroomObserver}, sourceType) {
		return false
	}

	return true
}

func (Ω DataSourceTargetID) ToInt() (int, error) {
	i, err := strconv.Atoi(string(Ω))
	if err != nil {
		return 0, errors.Wrapf(err, "Could not cast TargetID [%s] as int", Ω)
	}
	return i, nil
}

func NewDataSourceTargetIDFromInt(i int) (DataSourceTargetID, error) {
	if i == 0 {
		return DataSourceTargetID(""), errors.New("Invalid DataSourceTargetID: Received zero.")
	}
	return DataSourceTargetID(strconv.Itoa(i)), nil
}

type DataSourceTargetIDs []DataSourceTargetID

func (Ω DataSourceTargetIDs) Strings() (res []string) {
	for _, id := range Ω {
		res = append(res, string(id))
	}
	return
}

func (Ω DataSourceTargetIDs) AddToSet(ids ...DataSourceTargetID) DataSourceTargetIDs {
	for _, id := range ids {
		if Ω.Contains(id) {
			continue
		}
		Ω = append(Ω, id)
	}
	return Ω
}

func (Ω DataSourceTargetIDs) Contains(id DataSourceTargetID) bool {
	for i := range Ω {
		if Ω[i] == id {
			return true
		}
	}
	return false
}

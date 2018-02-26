package aggregate

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
)

func (Ω *Aggregate) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	res := map[nameusage.NameUsageID]nameusage.NameUsage{}
	for _, usage := range Ω.list {
		id, err := usage.ID()
		if err != nil {
			return nil, err
		}
		res[id] = usage
	}
	return json.Marshal(res)
}

func (Ω *Aggregate) UnmarshalJSON(provided []byte) error {

	Ω.list = []nameusage.NameUsage{}

	m := map[nameusage.NameUsageID]interface{}{}
	if err := json.Unmarshal(provided, &m); err != nil {
		return err
	}

	for id, i := range m {

		b, err := json.Marshal(i)
		if err != nil {
			return errors.Wrap(err, "Could not Marshal Usage interface")
		}

		usage, err := nameusage.NameUsageFromJSON(id, b)
		if err != nil {
			return err
		}
		Ω.list = append(Ω.list, usage)
	}

	return nil
}

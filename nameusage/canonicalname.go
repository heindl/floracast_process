package nameusage

import (
	"github.com/dropbox/godropbox/errors"
	"strings"
	"encoding/json"
	"bitbucket.org/heindl/taxa/utils"
)

// The canonical name is the scientific name of the species, subspecies, variety, etc. Anything under Genus.
type CanonicalName struct {
	name string
	rank string
}

func (a *CanonicalName) MarshalJSON() ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(map[string]interface{}{
		"ScientificName": a.name,
		"Rank": a.rank,
	})
}

func NewCanonicalName(name string, rank string) (*CanonicalName, error) {

	s := strings.TrimSpace(strings.ToLower(name))

	if s == "" {
		return nil, errors.New("Invalid CanonicalName: Received empty string")
	}

	wc := len(strings.Fields(s))
	if wc == 1 {
		return nil, errors.Newf("Invalid CanonicalName: Has only one word [%s], which suggests it is not a species or below", s)
	}
	if wc > 4 {
		return nil, errors.Newf("Invalid CanonicalName: More than four words [%s]", s)
	}

	// TODO: Consider removing some words like ".var" to improve consistency and matching.
	return &CanonicalName{
		name: s,
		rank: rank,
		}, nil
}

func (a *CanonicalName) Equals(b *CanonicalName) bool {
	return a.name == b.name
}

type CanonicalNames []*CanonicalName

func (Ω CanonicalNames) Strings() []string {
	res := []string{}
	for _, c := range Ω {
		res = utils.AddStringToSet(res, c.name)
	}
	return res
}

func (a CanonicalNames) Contains(b *CanonicalName) bool {
	for _, cn := range a {
		if cn.name == b.name {
			return true
		}
	}
	return false
}

func (Ω CanonicalNames) AddToSet(names ...*CanonicalName) CanonicalNames {
	for _, name := range names {
		if name == nil {
			continue
		}
		if !Ω.Contains(name) {
			Ω = append(Ω, name)
		}
	}
	return Ω
}
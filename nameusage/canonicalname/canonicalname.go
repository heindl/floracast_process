package canonicalname

import (
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

// The canonical Name is the scientific Name of the species, subspecies, variety, etc. Anything under Genus.
type CanonicalName struct {
	Name string `json:"ScientificName" firestore:"ScientificName"`
	Rank string `json:"Rank,omitempty" firestore:"Rank,omitempty"`
}

func (a *CanonicalName) ScientificName() string {
	if a == nil {
		return ""
	}
	return a.Name
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
	if wc > 5 {
		return nil, errors.Newf("Invalid CanonicalName: More than four words [%s]", s)
	}

	// TODO: Consider removing some words like ".var" to improve consistency and matching.
	return &CanonicalName{
		Name: s,
		Rank: rank,
	}, nil
}

func (a *CanonicalName) Equals(b *CanonicalName) bool {
	return a.Name == b.Name
}

type CanonicalNames []*CanonicalName

func (Ω CanonicalNames) ScientificNames() []string {
	res := []string{}
	for _, c := range Ω {
		res = utils.AddStringToSet(res, c.Name)
	}
	return res
}

func (a CanonicalNames) Contains(b ...*CanonicalName) bool {
	return utils.IntersectsStrings(a.ScientificNames(), CanonicalNames(b).ScientificNames())
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

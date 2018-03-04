package canonicalname

import (
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

// Name is the scientific name of the species, subspecies, variety, etc. Anything under genus.
type Name struct {
	SciName string `json:"ScientificName" firestore:"ScientificName"`
	Rank    string `json:"Rank,omitempty" firestore:"Rank,omitempty"`
}

// ScientificName returns just that.
func (Ω *Name) ScientificName() string {
	if Ω == nil {
		return ""
	}
	return Ω.SciName
}

// NewCanonicalName creates and validates a new name.
func NewCanonicalName(name string, rank string) (*Name, error) {

	s := strings.TrimSpace(strings.ToLower(name))

	if s == "" {
		return nil, errors.New("Invalid Name: Received empty string")
	}

	wc := len(strings.Fields(s))
	if wc == 1 {
		return nil, errors.Newf("Invalid Name: Has only one word [%s], which suggests it is not a species or below", s)
	}
	if wc > 5 {
		return nil, errors.Newf("Invalid Name: More than four words [%s]", s)
	}

	// TODO: Consider removing some words like ".var" to improve consistency and matching.
	return &Name{
		SciName: s,
		Rank:    rank,
	}, nil
}

//type TaxonRank string
//
//const (
//	// Originating from INaturalist:
//	RankKingdom     = TaxonRank("Kingdom")
//	RankPhylum      = TaxonRank("Phylum")
//	RankSubPhylum   = TaxonRank("SubPhylum")
//	RankClass       = TaxonRank("Class")
//	RankSubClass    = TaxonRank("SubClass")
//	RankOrder       = TaxonRank("Order")
//	RankSuperFamily = TaxonRank("SuperFamily")
//	RankFamily      = TaxonRank("Family")
//	RankSubFamily   = TaxonRank("SubFamily")
//	RankTribe       = TaxonRank("Tribe")
//	RankSubTribe    = TaxonRank("SubTribe")
//	RankGenus       = TaxonRank("Genus")
//	RankSpecies     = TaxonRank("Species")
//	RankSubSpecies  = TaxonRank("SubSpecies")
//	RankForm        = TaxonRank("Form")
//	RankVariety     = TaxonRank("Variety")
//)

// Equals is a helper function that could be flushed out to include Rank, but doesn't.
func (Ω *Name) Equals(b *Name) bool {
	return Ω.SciName == b.SciName
}

// Names is a utility tool for handling multiple names.
type Names []*Name

// ScientificNames returns all names.
func (Ω Names) ScientificNames() []string {
	res := []string{}
	for _, c := range Ω {
		res = utils.AddStringToSet(res, c.SciName)
	}
	return res
}

// Contains returns true if the Name is present in the list.
func (Ω Names) Contains(b ...*Name) bool {
	return utils.IntersectsStrings(Ω.ScientificNames(), Names(b).ScientificNames())
}

// AddToSet adds names to the list.
func (Ω Names) AddToSet(names ...*Name) Names {
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

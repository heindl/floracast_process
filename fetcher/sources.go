package main

import (
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/gbif"
)

func gatherSubspecies(name species.CanonicalName) ([]species.Species, error) {

	subspecies, err := gbif.Search(gbif.SearchQuery{
		Q:    string(name),
		Rank: "SUBSPECIES",
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not search gbif")
	}

	m := make(map[string][]int)

	addtoset := func(name string, n gbif.NameUsage) error {
		if n.Key == 0 {
			return errors.Wrapf(err, "no key found for subspecies:\n %v", utils.JsonOrSpew(n))
		}
		if _, ok := m[name]; ok {
			m[name] = utils.AddIntToSet(m[name], n.Key)
		} else {
			m[name] = []int{n.Key}
		}
		return nil
	}

	for _, sub := range subspecies {

		if err := addtoset(sub.CanonicalName, sub); err != nil {
			return nil, err
		}
		s := gbif.Species(sub.Key)
		synonyms, err := s.Synonyms()
		if err != nil {
			return nil, err
		}
		for _, synonym := range synonyms {
			if err := addtoset(sub.CanonicalName, synonym); err != nil {
				return nil, err
			}
		}
	}

	var response []species.Species

	for k, v := range m {
		s := species.Species{
			CanonicalName: species.CanonicalName(k),
		}
		for _, i := range v {
			s.Sources = append(s.Sources, species.Source{
				Type:     species.SourceTypeGBIF,
				IndexKey: species.IndexKey(i),
			})
		}
		response = append(response, s)
	}

	return response, nil

}

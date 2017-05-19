package main

import (
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/species/store"
	"bitbucket.org/heindl/utils"
	"github.com/heindl/gbif"
	"github.com/saleswise/errors/errors"
	"time"
	"github.com/heindl/eol"
	"flag"
	"gopkg.in/tomb.v2"
	. "github.com/saleswise/malias"
)

func main() {

	name := flag.String("name", "species canonical name", "species canonical name that will subdivided into lesser species")
	flag.Parse()


	store, err := store.NewSpeciesStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// http://localhost:4151/topic/create?topic=fetch-species
	fetcher := SpeciesFetcher{
		SpeciesStore: store,
	}

	if err := fetcher.FetchSpecies(species.CanonicalName(*name)); err != nil {
		panic(err)
	}
}

type SpeciesFetcher struct {
	SpeciesStore store.SpeciesStore
}

func (this *SpeciesFetcher) FetchSpecies(name species.CanonicalName) error {

	if !name.Valid() {
		return errors.New("canonical name not valid")
	}

	spcs, err := gatherSpecies(name)
	if err != nil {
		return err
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _spc := range spcs {
			s := _spc
			tmb.Go(func() error {
				for k, _ := range s.Sources {
					if err := this.SpeciesStore.AddSource(s.CanonicalName, k); err != nil {
						return err
					}
				}
				if err := this.SpeciesStore.SetClassification(s.CanonicalName, s.Classification); err != nil {
					return err
				}
				return this.setEOLMeta(s.CanonicalName)
			})
		}
		return nil
	})
	return tmb.Wait()
}

func gatherSpecies(name species.CanonicalName) ([]species.Species, error) {

	list := make(map[species.CanonicalName]species.Species)

	subspecies, err := gbif.Search(gbif.SearchQuery{
		Q:    string(name),
		Rank: []gbif.Rank{gbif.RankSUBSPECIES, gbif.RankSPECIES},
		Status: []gbif.TaxonomicStatus{gbif.TaxonomicStatusACCEPTED},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not search gbif")
	}

	for _, sp := range subspecies {
		if sp.TaxonomicStatus != gbif.TaxonomicStatusACCEPTED {
			continue
		}
		var s species.Species
		if sp.NubKey == 0 {
			continue
		}
		if _, ok := list[species.CanonicalName(sp.CanonicalName)]; !ok {
			s = species.Species{
				CanonicalName: species.CanonicalName(sp.CanonicalName),
				ScientificName: sp.ScientificName,
				CreatedAt: utils.TimePtr(time.Now()),
				Sources: map[species.SourceKey]species.SourceData{
					species.NewSourceKey(species.IndexKey(sp.NubKey), species.SourceTypeGBIF): species.SourceData{},
				},
				Classification: sp.Classification,
			}
		} else {
			s = list[species.CanonicalName(sp.CanonicalName)]
		}
		synonyms, err := gbif.Species(sp.Key).Synonyms()
		if err != nil {
			return nil, err
		}
		for _, sy := range synonyms {
			//if sy.TaxonomicStatus == gbif.TaxonomicStatusACCEPTED {
			//	s.Classification = sy.Classification
			//}
			if sy.NubKey == 0 {
				continue
			}
			k := species.NewSourceKey(species.IndexKey(sy.NubKey), species.SourceTypeGBIF)
			if _, ok := s.Sources[k]; ok {
				continue
			}
			s.Sources[k] = species.SourceData{}
		}
		list[species.CanonicalName(sp.CanonicalName)] = s
	}


	var response []species.Species
	for _, v := range list {
		response = append(response, v)
	}

	return response, nil

}

func (this *SpeciesFetcher) setEOLMeta(name species.CanonicalName) error {

	if !name.Valid() {
		return errors.New("canonical name not valid")
	}

	results, err := eol.Search(eol.SearchQuery{
		Query: string(name),
		Limit: 10,
	})
	if err != nil {
		return errors.Wrap(err, "could not search encyclopedia of life").SetState(M{utils.LogkeyCanonicalName: name})
	}

	if len(results) == 0 {
		return nil
	}

	// The first result should be the most relevant, but check the top ten for the highest score.

	var highest eol.PageResponse

	for _, r := range results {

		page, err := eol.Page(eol.PageQuery{
			ID:      r.ID,
			Images:  1,
			Text:    1,
			Details: true,
		})
		if err != nil {
			return errors.Wrapf(err, "could not find page query from id[%v]", r.ID)
		}
		if page.RichnessScore > highest.RichnessScore {
			highest = *page
		}
	}

	if highest.Identifier == 0 {
		return nil
	}

	if err := this.SpeciesStore.AddSource(species.CanonicalName(name), species.NewSourceKey(species.IndexKey(highest.Identifier), species.SourceTypeEOL)); err != nil {
		return err
	}

	if len(highest.Texts()) > 0 {
		if err := this.SpeciesStore.SetDescription(name, &species.Media{
			Source: highest.Texts()[0].Source,
			Value:  highest.Texts()[0].Value,
		}); err != nil {
			return err
		}
	}

	if len(highest.Images()) > 0 {
		if err := this.SpeciesStore.SetImage(name, &species.Media{
			Source: highest.Images()[0].Source,
			Value:  highest.Images()[0].Value,
		}); err != nil {
			return err
		}
	}

	return nil
}

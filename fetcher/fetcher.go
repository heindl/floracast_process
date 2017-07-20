package main
//
//import (
//	"bitbucket.org/heindl/species"
//	"bitbucket.org/heindl/species/store"
//	"bitbucket.org/heindl/utils"
//	"github.com/heindl/gbif"
//	"github.com/saleswise/errors/errors"
//	"time"
//	"github.com/heindl/eol"
//	"flag"
//	"fmt"
//	"strings"
//)
//
//func main() {
//
//	name := flag.String("name", "species canonical name", "species canonical name that will subdivided into lesser species")
//	flag.Parse()
//
//	store, err := store.NewSpeciesStore()
//	if err != nil {
//		panic(err)
//	}
//	defer store.Close()
//
//	// http://localhost:4151/topic/create?topic=fetch-species
//	fetcher := SpeciesFetcher{
//		SpeciesStore: store,
//	}
//
//	if err := fetcher.FetchSpecies(store.CanonicalName(*name)); err != nil {
//		panic(err)
//	}
//}
//
//type SpeciesFetcher struct {
//	SpeciesStore store.TaxaStore
//}
//
//
//func (Ω *SpeciesFetcher) FetchSpecies(name store.CanonicalName) error {
//
//	if !name.Valid() {
//		return errors.New("canonical name not valid")
//	}
//
//	subspecies, err := gbif.Search(gbif.SearchQuery{
//		Q:    string(name),
//		Rank: []gbif.Rank{gbif.RankSUBSPECIES, gbif.RankSPECIES},
//		Status: []gbif.TaxonomicStatus{gbif.TaxonomicStatusACCEPTED},
//	})
//	if err != nil {
//		return nil, errors.Wrap(err, "could not search gbif")
//	}
//
//	for _, sp := range subspecies {
//
//		if sp.TaxonomicStatus != gbif.TaxonomicStatusACCEPTED {
//			continue
//		}
//		if sp.NubKey == 0 {
//			continue
//		}
//		if strings.TrimSpace(sp.CanonicalName) == "" {
//			continue
//		}
//
//		// Ensure the string begins with the initial value of the string. 'Morchella' is namely a problem case.
//		if !strings.HasPrefix(sp.CanonicalName, strings.ToLower(string(name))) {
//			continue
//		}
//
//		if err := Ω.FetchSubspecies(store.CanonicalName(sp.CanonicalName)); err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//
//}
//
//func (Ω *SpeciesFetcher) FetchSubspecies(usage gbif.NameUsage) error {
//
//	key := species.NewSpeciesKey(&store.CanonicalName(usage.CanonicalName))
//
//	if err := Ω.SpeciesStore.SetClassification(key, usage.Classification); err != nil {
//		return err
//	}
//
//	var s *species.Species
//	if i := specieslist.Index(species.NewSpeciesKey(subspeciesName)); i == -1 {
//		specieslist = append(specieslist, species.Species{
//			Key: species.NewSpeciesKey(subspeciesName),
//			ScientificNames: []string{sp.ScientificName},
//			CreatedAt: utils.TimePtr(time.Now()),
//			Classification: sp.Classification,
//		})
//	}
//	s = specieslist[specieslist.Index(species.NewSpeciesKey(subspeciesName))]
//	s.ScientificNames = utils.AddStringToSet(s.ScientificNames, sp.ScientificName)
//	sourcelist = sourcelist.AddToSet(species.NewDataSourceKey(
//		species.IndexKey(sp.NubKey),
//		species.SourceTypeGBIF,
//		subspeciesName,
//	))
//
//	synonyms, err := gbif.Species(sp.Key).Synonyms()
//	if err != nil {
//		return nil, err
//	}
//	for _, sy := range synonyms {
//		//if sy.TaxonomicStatus == gbif.TaxonomicStatusACCEPTED {
//		//	s.Classification = sy.Classification
//		//}
//		if sy.NubKey == 0 {
//			continue
//		}
//		sourcelist = sourcelist.AddToSet(species.NewDataSourceKey(
//			species.IndexKey(sy.NubKey),
//			species.SourceTypeGBIF,
//			name,
//		))
//		s.ScientificNames = utils.AddStringToSet(s.ScientificNames, sy.ScientificName)
//	}
//
//}
//
//func (this *SpeciesFetcher) setEOLMeta(name species.CanonicalName) error {
//
//	if !name.Valid() {
//		return errors.Newf("canonical name not valid: %s", name)
//	}
//
//	results, err := eol.Search(eol.SearchQuery{
//		Query: string(name),
//		Limit: 10,
//	})
//	if err != nil {
//		return errors.Wrap(err, "could not search encyclopedia of life").SetState(M{utils.LogkeyCanonicalName: name})
//	}
//
//	if len(results) == 0 {
//		return nil
//	}
//
//	// The first result should be the most relevant, but check the top ten for the highest score.
//
//	var highest eol.PageResponse
//
//	for _, r := range results {
//
//		page, err := eol.Page(eol.PageQuery{
//			ID:      r.ID,
//			Images:  1,
//			Text:    1,
//			Details: true,
//			CommonNames: true,
//
//		})
//		if err != nil && !eol.ErrNotFound(err) {
//			return errors.Wrapf(err, "could not find page query from id[%v]", r.ID)
//		}
//		if err != nil && eol.ErrNotFound(err) {
//			// Log error message here
//			fmt.Println("Could not find EOL page: ", r.ID)
//			return nil
//		}
//		if page.RichnessScore > highest.RichnessScore {
//			highest = *page
//		}
//	}
//
//	if highest.Identifier == 0 {
//		return nil
//	}
//
//	for _, vn := range highest.VernacularNames {
//		if vn.Language == "en" && vn.EolPreferred {
//			if err := this.SpeciesStore.SetCommonName(name, vn.VernacularName); err != nil {
//				return err
//			}
//			break
//		}
//	}
//
//	if err := this.SpeciesStore.AddSource(species.NewDataSourceKey(species.IndexKey(highest.Identifier), species.SourceTypeEOL, name)); err != nil {
//		return err
//	}
//
//	if len(highest.Texts()) > 0 {
//		var media []*species.Media
//		for _, t := range highest.Texts() {
//			media = append(media, &species.Media{
//				Source: t.Source,
//				Value: t.Description,
//			})
//		}
//		if err := this.SpeciesStore.SetDescriptions(name, media...); err != nil {
//			return err
//		}
//	}
//
//	if len(highest.Images()) > 0 {
//		var media []*species.Media
//		for _, t := range highest.Images() {
//			media = append(media, &species.Media{
//				Source: t.Source,
//				Value: t.MediaURL,
//			})
//		}
//		if err := this.SpeciesStore.SetImages(name, media...); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

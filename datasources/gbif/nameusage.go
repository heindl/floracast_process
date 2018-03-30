package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"gopkg.in/tomb.v2"
	"strings"
	"sync"
)

type orchestrator struct {
	Usages []nameusage.NameUsage
	sync.Mutex
}

func (Ω *orchestrator) hasCanonicalName(name string) (bool, error) {
	Ω.Lock()
	defer Ω.Unlock()
	for _, u := range Ω.Usages {
		hasName, err := u.HasScientificName(name)
		if err != nil {
			return false, err
		}
		if hasName {
			return true, nil
		}
	}
	return false, nil
}

// FetchNameUsages matches the canonical name and taxon ids to generate NameUsages.
func FetchNameUsages(cxt context.Context, namesToMatch []string, keysToMatch datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	o := orchestrator{
		Usages: []nameusage.NameUsage{},
	}

	// First match keys.
	if err := o.matchKeys(keysToMatch); err != nil {
		return nil, err
	}

	// Recursively fetch them all
	if err := o.matchNames(namesToMatch); err != nil {
		return nil, err
	}

	return o.Usages, nil
}

func (Ω *orchestrator) matchKeys(keysToMatch datasources.TargetIDs) error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _key := range keysToMatch {
			key := _key
			tmb.Go(func() error {
				k, err := targetIDToSpecies(key)
				if err != nil {
					return err
				}
				if err := Ω.matchKey(k); err != nil && err != errUnsupported {
					return err
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *orchestrator) matchNames(namesToMatch []string) error {

	names := utils.StringsToLower(utils.RemoveStringDuplicates(namesToMatch)...)

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {
				// Run Name to check for error
				canonicalName, err := canonicalname.NewCanonicalName(name, "")
				if err != nil {
					return err
				}
				hasName, err := Ω.hasCanonicalName(canonicalName.ScientificName())
				if err != nil {
					return err
				}
				if hasName {
					return nil
				}

				res, err := match(matchQuery{Name: name})
				if err != nil {
					return err
				}

				if !res.UsageKey.Valid() {
					fmt.Println(fmt.Sprintf("Canonical name [%s] not matched in GBIF.", name))
				}

				return Ω.matchKey(res.UsageKey)
			})
		}
		return nil
	})
	return tmb.Wait()
}

var errUnsupported = fmt.Errorf("unsupported usage")

func (Ω *orchestrator) matchKey(usageKey species) error {
	// Get the reference for the synonym.
	//gbifNameUsage := nameUsage{}

	//if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d", usageKey), &gbifNameUsage); err != nil {
	//	return err
	//}

	gbifUsage, err := species(int(usageKey)).Name()
	if err != nil {
		return err
	}

	if gbifUsage.Rank == rankGenus {
		fmt.Println(fmt.Sprintf("Warning: Encountered genus [%d], but unsupported", gbifUsage.Key))
		return errUnsupported
	}

	parentKey := gbifUsage.AcceptedKey
	key := gbifUsage.Key

	if gbifUsage.Synonym {
		return Ω.matchKey(parentKey)
	}
	return Ω.parseNameUsage(gbifUsage.CanonicalName, string(gbifUsage.Rank), key)

}

func (Ω *orchestrator) parseSource(scientificName, rank string, spcs species) (nameusage.Source, error) {
	canonicalName, err := canonicalname.NewCanonicalName(scientificName, strings.ToLower(rank))
	if err != nil {
		return nil, err
	}

	usageSource, err := nameusage.NewSource(datasources.TypeGBIF, spcs.TargetID(), canonicalName)
	if err != nil {
		fmt.Println(fmt.Sprintf("WARNING: Invalid GBIF nameUsage [%s, %d]", canonicalName.ScientificName(), spcs))
		return nil, nil
	}

	vernacularNames, err := species(int(spcs)).fetchVernacularNames()
	if err != nil {
		return nil, err
	}
	for _, vn := range vernacularNames {
		if vn.Language == "eng" || vn.Language == "" {
			if err = usageSource.AddCommonNames(strings.TrimSpace(vn.VernacularName)); err != nil {
				return nil, err
			}
		}
	}
	return usageSource, nil
}

func (Ω *orchestrator) parseNameUsage(scientificName, rank string, spcs species) error {

	usageSource, err := Ω.parseSource(scientificName, rank, spcs)
	if err != nil {
		return err
	}
	// Suggests an invalid source.
	if usageSource == nil {
		return nil
	}

	usage, err := nameusage.NewNameUsage(usageSource)
	if err != nil {
		return err
	}

	synonymUsageSources, err := matchSynonyms(spcs)
	if err != nil {
		return err
	}

	if err := usage.AddSources(synonymUsageSources...); err != nil {
		return err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.Usages = append(Ω.Usages, usage)

	return nil
}

func matchSynonyms(id species) ([]nameusage.Source, error) {

	synonymUsages, err := id.Synonyms()
	if err != nil {
		return nil, err
	}

	response := []nameusage.Source{}

	for _, synonym := range synonymUsages {
		src, err := parseSynonymsAsSource(synonym)
		if err != nil {
			return nil, err
		}
		if src == nil {
			continue
		}
		response = append(response, src)
	}

	return response, nil
}

func parseSynonymsAsSource(synonym *nameUsage) (nameusage.Source, error) {
	if synonym.Rank == rankUnranked {
		fmt.Println(fmt.Sprintf("Warning: usage [%d] is unranked, skipping", synonym.Key))
		return nil, nil
	}

	acceptedTaxonomicStatuses := taxonomicStatuses{
		taxonomicStatusSynonym,
		taxonomicStatusHeterotypicSynonym,
		taxonomicStatusHomotypicSynonym,
		taxonomicStatusProparteSynonym,
	}

	if !acceptedTaxonomicStatuses.Contains(synonym.TaxonomicStatus) || !synonym.Synonym {
		fmt.Println(fmt.Sprintf("Warning: usage [%d] is not a synonym [%s], so skipping", synonym.Key, synonym.TaxonomicStatus))
		return nil, nil
	}

	if strings.IndexFunc(synonym.CanonicalName, func(r rune) bool {
		return (r < 'A' || r > 'z') && r != ' ' && r != '-'
	}) != -1 {
		fmt.Println(fmt.Sprintf("Warning: name [%s] contains non letter, so skipping", synonym.CanonicalName))
		return nil, nil
	}

	canonicalName, err := canonicalname.NewCanonicalName(synonym.CanonicalName, strings.ToLower(string(synonym.Rank)))
	if err != nil {
		return nil, err
	}

	src, err := nameusage.NewSource(datasources.TypeGBIF, synonym.Key.TargetID(), canonicalName)
	if err != nil {
		return nil, err
	}

	return src, nil
}

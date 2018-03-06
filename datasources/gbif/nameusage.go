package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"net/url"
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
				k, err := taxonIDFromTargetID(key)
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
				return Ω.matchName(name)
			})
		}
		return nil
	})
	return tmb.Wait()
}

func taxonIDFromTargetID(id datasources.TargetID) (taxonID, error) {
	i, err := id.ToInt()
	if err != nil {
		return taxonID(0), errors.Wrapf(err, "Could not cast GBIF TargetID [%s] as taxonID", id)
	}
	txnID := taxonID(i)
	if !txnID.Valid() {
		return taxonID(0), errors.Newf("Invalid GBIF taxonID [%s]", id)
	}
	return txnID, nil
}

func (Ω *orchestrator) matchName(name string) error {

	queryName := utils.CapitalizeString(url.QueryEscape(strings.ToLower(name)))

	u := fmt.Sprintf("http://api.gbif.org/v1/species/match?name=%s&verbose=true", queryName)
	matchResult := matchResult{}
	if err := utils.RequestJSON(u, &matchResult); err != nil {
		return err
	}

	if !matchResult.UsageKey.Valid() {
		fmt.Println(fmt.Sprintf("Canonical name [%s] not matched in GBIF.", name))
	}

	//if matchResult.rank == api.rankGenus {
	//	return errors.Newf("Unexpected GENUS [%d] returned in species match", matchResult.UsageKey)
	//}

	// We will fetch synonyms later, so just match the key
	//if matchResult.Synonym {
	// We need the vernacular name, and other information, so fetch key.
	return Ω.matchKey(matchResult.UsageKey)
	//}

	//return Ω.parseNameUsage(matchResult.Name, matchResult.vernacularName, string(matchResult.rank), matchResult.UsageKey)
}

var errUnsupported = fmt.Errorf("unsupported usage")

func (Ω *orchestrator) matchKey(usageKey taxonID) error {
	// Get the reference for the synonym.
	gbifNameUsage := nameUsage{}

	if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d", usageKey), &gbifNameUsage); err != nil {
		return err
	}

	if gbifNameUsage.Rank == rankGenus {
		fmt.Println(fmt.Sprintf("Warning: Encountered genus [%d], but unsupported", gbifNameUsage.Key))
		return errUnsupported
	}

	parentKey := gbifNameUsage.AcceptedKey
	key := gbifNameUsage.Key

	if gbifNameUsage.Synonym {
		return Ω.matchKey(parentKey)
	}
	return Ω.parseNameUsage(gbifNameUsage.CanonicalName, gbifNameUsage.VernacularName, string(gbifNameUsage.Rank), key)

}

func (Ω *orchestrator) parseNameUsage(scientificName, vernacularName, rank string, taxonID taxonID) error {
	canonicalName, err := canonicalname.NewCanonicalName(scientificName, strings.ToLower(rank))
	if err != nil {
		return nil
	}

	usageSource, err := nameusage.NewSource(datasources.TypeGBIF, taxonID.TargetID(), canonicalName)
	if err != nil {
		fmt.Println(fmt.Sprintf("WARNING: Invalid GBIF nameUsage [%s, %d]", canonicalName.ScientificName(), taxonID))
		return nil
	}

	vernacularNames, err := species(int(taxonID)).fetchVernacularNames()
	if err != nil {
		return err
	}
	for _, vn := range vernacularNames {
		if vn.Language != "eng" && vn.Language != "" {
			continue
		}
		if err = usageSource.AddCommonNames(strings.TrimSpace(vn.VernacularName)); err != nil {
			return err
		}
	}

	usage, err := nameusage.NewNameUsage(usageSource)
	if err != nil {
		return err
	}

	synonymUsageSources, err := matchSynonyms(taxonID)
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

func matchSynonyms(id taxonID) ([]nameusage.Source, error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf("http://api.gbif.org/v1/species/%d/synonyms?", id))
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

func fetchNameUsages(url string) ([]*nameUsage, error) {

	offset := 0
	records := []*nameUsage{}

	for {
		var res struct {
			Offset       int          `json:"offset"`
			Limit        int          `json:"limit"`
			EndOfRecords bool         `json:"endOfRecords"`
			Results      []*nameUsage `json:"results"`
		}

		if err := utils.RequestJSON(
			url+fmt.Sprintf("&offset=%d&limit=300", offset),
			&res,
		); err != nil {
			return nil, err
		}

		records = append(records, res.Results...)
		if res.EndOfRecords {
			break
		}
		offset++
	}

	return records, nil
}

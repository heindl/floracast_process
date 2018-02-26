package gbif

import (
	"context"
	"fmt"
	"net/url"
	"bitbucket.org/heindl/process/utils"
	"gopkg.in/tomb.v2"
	"strings"
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif/api"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"sync"
)

type orchestrator struct {
	Usages []nameusage.NameUsage
	Context context.Context
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

func FetchNamesUsages(cxt context.Context, namesToMatch []string, keysToMatch datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	o := orchestrator{
		Usages: []nameusage.NameUsage{},
		Context: cxt,
	}

	// First match keys.
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _key := range keysToMatch {
			key := _key
			tmb.Go(func() error {
				k, err := TaxonIDFromTargetID(key)
				if err != nil {
					return err
				}
				if err := o.matchKey(k); err != nil && err != ErrUnsupported {
					return err
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	// Recursively fetch them all
	tmb = tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range utils.StringsToLower(utils.RemoveStringDuplicates(namesToMatch)...) {
			name := _name
			tmb.Go(func() error {
				// Run CanonicalName to check for error
				canonicalName, err := canonicalname.NewCanonicalName(name, "")
				if err != nil {
					return err
				}
				hasName, err := o.hasCanonicalName(canonicalName.ScientificName())
				if err != nil {
					return err
				}
				if hasName {
					return nil
				}
				return o.matchName(name)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return o.Usages, nil
}

func TaxonIDFromTargetID(id datasources.TargetID) (api.TaxonID, error) {
	i, err := id.ToInt()
	if err != nil {
		return api.TaxonID(0), errors.Wrapf(err, "Could not cast GBIF TargetID [%s] as TaxonID", id)
	}
	taxonID := api.TaxonID(i)
	if !taxonID.Valid() {
		return api.TaxonID(0), errors.Newf("Invalid GBIF TaxonID [%s]", id)
	}
	return taxonID, nil
}

func (Ω *orchestrator) matchName(name string) error {

	u := fmt.Sprintf("http://api.gbif.org/v1/species/match?name=%s&verbose=true", url.QueryEscape(name))

	matchResult := MatchResult{}
	if err := utils.RequestJSON(u, &matchResult); err != nil {
		return err
	}

	if !matchResult.UsageKey.Valid() {
		fmt.Println(fmt.Sprintf("Canonical name [%s] not matched in GBIF.", name))
	}

	//if matchResult.Rank == api.RankGENUS {
	//	return errors.Newf("Unexpected GENUS [%d] returned in species match", matchResult.UsageKey)
	//}

	// We will fetch synonyms later, so just match the key
	//if matchResult.Synonym {
	// We need the vernacular name, and other information, so fetch key.
	return Ω.matchKey(matchResult.UsageKey)
	//}

	//return Ω.fashionCanonicalNameUsage(matchResult.CanonicalName, matchResult.VernacularName, string(matchResult.Rank), matchResult.UsageKey)
}

var ErrUnsupported = fmt.Errorf("unsupported usage")

func (Ω *orchestrator) matchKey(usageKey api.TaxonID) error {
	// Get the reference for the synonym.
	gbifNameUsage := api.NameUsage{}
	if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d?", usageKey), &gbifNameUsage); err != nil {
		return err
	}

	if gbifNameUsage.Rank == api.RankGENUS {
		fmt.Println(fmt.Sprintf("Warning: Encountered genus [%d], but unsupported", gbifNameUsage.Key))
		return ErrUnsupported
	}

	parentKey := gbifNameUsage.AcceptedKey
	key := gbifNameUsage.Key

	if gbifNameUsage.Synonym {
		return Ω.matchKey(parentKey)
	}

	return Ω.fashionCanonicalNameUsage(gbifNameUsage.CanonicalName, gbifNameUsage.VernacularName, string(gbifNameUsage.Rank), key)

}

func (Ω *orchestrator) fashionCanonicalNameUsage(scientificName, vernacularName, rank string, taxonID api.TaxonID) error {
	canonicalName, err := canonicalname.NewCanonicalName(scientificName, strings.ToLower(rank))
	if err != nil {
		return nil
	}

	usageSource, err := nameusage.NewSource(datasources.TypeGBIF, taxonID.TargetID(), canonicalName)
	if err != nil {
		fmt.Println(fmt.Sprintf("WARNING: Invalid GBIF NameUsage [%s, %d]", canonicalName.ScientificName(), taxonID))
		return nil
	}

	if vernacularName != "" {
		if err := usageSource.AddCommonNames(vernacularName); err != nil {
			return nil
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


type MatchResult struct {
	UsageKey       api.TaxonID    `json:"usageKey"`
	ScientificName string `json:"scientificName"`
	CanonicalName  string `json:"canonicalName"`
	Rank           api.Rank `json:"rank"`
	Status         string `json:"status"`
	Confidence     int    `json:"confidence"`
	Note           string `json:"note"`
	MatchType      string `json:"matchType"`
	Synonym    bool   `json:"synonym"`
}

func matchSynonyms(id api.TaxonID) ([]nameusage.Source, error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf( "http://api.gbif.org/v1/species/%d/synonyms?", id))
	if err != nil {
		return nil, err
	}

	response := []nameusage.Source{}

	for _, synonym := range synonymUsages {

		if synonym.Rank == api.RankUNRANKED {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is unranked, skipping", synonym.Key))
			continue
		}

		acceptedTaxonomicStatuses := api.TaxonomicStatuses{
			api.TaxonomicStatusSYNONYM,
			api.TaxonomicStatusHETEROTYPIC_SYNONYM,
			api.TaxonomicStatusHOMOTYPIC_SYNONYM,
			api.TaxonomicStatusPROPARTE_SYNONYM,
			}

		if !acceptedTaxonomicStatuses.Contains(synonym.TaxonomicStatus) || !synonym.Synonym {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is not a synonym [%s], so skipping", synonym.Key, synonym.TaxonomicStatus))
			continue
		}

		if strings.IndexFunc(synonym.CanonicalName, func(r rune) bool {
			return (r < 'A' || r > 'z') && r != ' ' && r != '-'
		}) != -1 {
			fmt.Println(fmt.Sprintf("Warning: name [%s] contains non letter, so skipping", synonym.CanonicalName))
			continue
		}

		canonicalName, err := canonicalname.NewCanonicalName(synonym.CanonicalName, strings.ToLower(string(synonym.Rank)))
		if err != nil {
			return nil, err
		}

		src, err := nameusage.NewSource(datasources.TypeGBIF, synonym.Key.TargetID(), canonicalName)
		if err != nil {
			return nil, err
		}

		response = append(response, src)
	}

	return response, nil
}

func fetchNameUsages(url string) ([]*api.NameUsage, error) {

	offset := 0
	records := []*api.NameUsage{}

	for {
		var res struct {
			Offset int `json:"offset"`
			Limit int `json:"limit"`
			EndOfRecords bool `json:"endOfRecords"`
			Results []*api.NameUsage `json:"results"`
		}

		nUrl := url + fmt.Sprintf("&offset=%d&limit=300", offset)
		if err := utils.RequestJSON(nUrl, &res); err != nil {
			return nil, err
		}

		records = append(records, res.Results...)
		if res.EndOfRecords {
			break
		}
		offset += 1
	}

	return records, nil
}
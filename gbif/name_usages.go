package gbif

import (
	"context"
	"fmt"
	"net/url"
	"bitbucket.org/heindl/taxa/utils"
	"gopkg.in/tomb.v2"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/gbif/api"
	"strconv"
)

type orchestrator struct {
	Usages *name_usage.AggregateNameUsages
	Context context.Context
}

func FetchNamesUsages(cxt context.Context, namesToMatch []string, keysToMatch store.DataSourceTargetIDs) (*name_usage.AggregateNameUsages, error) {

	o := orchestrator{
		Usages: &name_usage.AggregateNameUsages{},
		Context: cxt,
	}

	// First match keys.
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _key := range keysToMatch {
			key := _key
			tmb.Go(func() error {
				if err := o.matchKey(TaxonIDFromTargetID(key)); err != nil && err != ErrUnsupported {
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
				canonicalName, err := name_usage.NewCanonicalName(name, "")
				if err != nil {
					return err
				}
				if i := o.Usages.FirstIndexOfName(canonicalName); i == -1 {
					return o.matchName(name)
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return o.Usages, nil
}

func TaxonIDFromTargetID(id store.DataSourceTargetID) api.TaxonID {
	i, err := strconv.Atoi(string(id))
	if err != nil {
		return api.TaxonID(0)
	}
	return api.TaxonID(i)
}

func (Ω *orchestrator) matchName(name string) error {

	u := fmt.Sprintf("http://api.gbif.org/v1/species/match?name=%s&verbose=true", url.QueryEscape(name))

	matchResult := MatchResult{}
	if err := utils.RequestJSON(u, &matchResult); err != nil {
		return err
	}

	if matchResult.Rank == api.RankGENUS {
		return errors.Newf("Unexpected GENUS [%d] returned in species match", matchResult.UsageKey)
	}

	// We will fetch synonyms later, so just match the key
	if matchResult.Synonym {
		return Ω.matchKey(matchResult.UsageKey)
	}

	canonicalName, err := name_usage.NewCanonicalName(matchResult.CanonicalName, strings.ToLower(string(matchResult.Rank)))
	if err != nil {
		return nil
	}

	usageSource, err := name_usage.NewNameUsageSource(store.DataSourceTypeGBIF, matchResult.UsageKey.TargetID(), canonicalName, true)
	if err != nil {
		return err
	}

	usage, err := name_usage.NewCanonicalNameUsage(usageSource)
	if err != nil {
		return err
	}

	synonymUsageSources, err := matchSynonyms(matchResult.UsageKey)
	if err != nil {
		return err
	}

	if err := usage.AddSynonyms(synonymUsageSources...); err != nil {
		return err
	}

	if err := Ω.Usages.Add(usage); err != nil {
		return err
	}

	return nil
}

var ErrUnsupported = fmt.Errorf("unsupported usage")

func (Ω *orchestrator) matchKey(usageKey api.TaxonID) error {
	// Get the reference for the synonym.
	nameUsage := api.NameUsage{}
	if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d?", usageKey), &nameUsage); err != nil {
		return err
	}

	if nameUsage.Rank == api.RankGENUS {
		fmt.Println(fmt.Sprintf("Warning: Encountered genus [%d], but unsupported", nameUsage.Key))
		return ErrUnsupported
	}

	parentKey := nameUsage.AcceptedKey
	key := nameUsage.Key

	if nameUsage.Synonym {
		return Ω.matchKey(parentKey)
	}

	canonicalName, err := name_usage.NewCanonicalName(nameUsage.CanonicalName, strings.ToLower(string(nameUsage.Rank)))
	if err != nil {
		return nil
	}

	usageSource, err := name_usage.NewNameUsageSource(store.DataSourceTypeGBIF, key.TargetID(), canonicalName, true)
	if err != nil {
		return err
	}

	usage, err := name_usage.NewCanonicalNameUsage(usageSource)
	if err != nil {
		return err
	}

	synonymUsageSources, err := matchSynonyms(key)
	if err != nil {
		return err
	}

	if err := usage.AddSynonyms(synonymUsageSources...); err != nil {
		return err
	}

	if err := Ω.Usages.Add(usage); err != nil {
		return err
	}

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

func matchSynonyms(id api.TaxonID) ([]*name_usage.NameUsageSource, error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf( "http://api.gbif.org/v1/species/%d/synonyms?", id))
	if err != nil {
		return nil, err
	}

	response := []*name_usage.NameUsageSource{}

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

		canonicalName, err := name_usage.NewCanonicalName(synonym.CanonicalName, strings.ToLower(string(synonym.Rank)))
		if err != nil {
			return nil, err
		}

		src, err := name_usage.NewNameUsageSource(store.DataSourceTypeGBIF, synonym.Key.TargetID(), canonicalName, true)
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
package gbif

import (
	"context"
	"fmt"
	"net/url"
	"bitbucket.org/heindl/taxa/utils"
	"gopkg.in/tomb.v2"
	"sync"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"bitbucket.org/heindl/taxa/store"
)

type orchestrator struct {
	Usages name_usage.CanonicalNameUsages
	OccurrenceCount map[TaxonID]int
	sync.Mutex
	Context context.Context
}

func FetchNamesUsages(cxt context.Context, namesToMatch []string, keysToMatch store.DataSourceTargetIDs) (name_usage.CanonicalNameUsages, error) {

	o := orchestrator{
		Usages: name_usage.CanonicalNameUsages{},
		OccurrenceCount: map[TaxonID]int{},
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


	// First let's get all the names together and unique.
	names := utils.RemoveStringDuplicates(namesToMatch)
	names = utils.StringsToLower(names...)

	// Recursively fetch them all
	tmb = tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {
				n := strings.ToLower(name)
				if i := o.Usages.FirstIndexOfName(n); i == -1 {
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

	var err error
	o.Usages, err = o.Usages.Condense()
	if err != nil {
		return nil, err
	}

	// Load cache with occurrence counts
	tmb = tomb.Tomb{}
	tmb.Go(func() error {
		for _, _id := range o.Usages.TargetIDs(store.DataSourceTypeGBIF) {
			id := _id
			tmb.Go(func() error {
				count, err := o.occurrenceCount(TaxonIDFromTargetID(id))
				if err != nil {
					return err
				}
				o.Lock()
				defer o.Unlock()
				for i := range o.Usages {
					if o.Usages[i].SourceTargetOccurrenceCount.Contains(store.DataSourceTypeGBIF, id) {
						o.Usages[i].SourceTargetOccurrenceCount.Set(store.DataSourceTypeGBIF, id, count)
					}
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

func (Ω *orchestrator) matchName(name string) error {

	u := fmt.Sprintf("http://api.gbif.org/v1/species/match?name=%s&verbose=true", url.QueryEscape(name))

	matchResult := MatchResult{}
	if err := utils.RequestJSON(u, &matchResult); err != nil {
		return err
	}

	if strings.ToLower(matchResult.Rank) == "genus" {
		return errors.Newf("Unexpected GENUS [%d] returned in species match", matchResult.UsageKey)
	}

	// We will fetch synonyms later, so just match the key
	if matchResult.Synonym {
		return Ω.matchKey(matchResult.UsageKey)
	}

	names, ranks, sourceTargetOccurrenceCount, err := Ω.matchSynonyms(matchResult.UsageKey)
	if err != nil {
		return err
	}

	canonicalNameUsage := name_usage.CanonicalNameUsage{
		CanonicalName: strings.ToLower(matchResult.CanonicalName),
		Synonyms: names,
		Ranks: utils.AddStringToSet(ranks, strings.ToLower(matchResult.Rank)),
		SourceTargetOccurrenceCount: sourceTargetOccurrenceCount,
	}
	canonicalNameUsage.SourceTargetOccurrenceCount.Set(store.DataSourceTypeGBIF, matchResult.UsageKey.TargetID(), 0)


	Ω.Lock()
	defer Ω.Unlock()

	Ω.Usages = append(Ω.Usages, canonicalNameUsage)

	return nil
}

var ErrUnsupported = fmt.Errorf("unsupported usage")

func (Ω *orchestrator) matchKey(usageKey TaxonID) error {
	// Get the reference for the synonym.
	nameUsage := NameUsage{}
	if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d?", usageKey), &nameUsage); err != nil {
		return err
	}

	if nameUsage.Rank == "GENUS" {
		fmt.Println(fmt.Sprintf("Warning: Encountered genus [%d], but unsupported", nameUsage.Key))
		return ErrUnsupported
	}

	parentKey := nameUsage.AcceptedKey
	key := nameUsage.Key

	if nameUsage.Synonym {
		return Ω.matchKey(parentKey)
	}

	names, ranks, mapTaxonIDCounts, err := Ω.matchSynonyms(key)
	if err != nil {
		return err
	}

	canonicalNameUsage := name_usage.CanonicalNameUsage{
		CanonicalName: strings.ToLower(nameUsage.CanonicalName),
		Synonyms: names,
		Ranks: utils.AddStringToSet(ranks, strings.ToLower(nameUsage.Rank)),
		SourceTargetOccurrenceCount: mapTaxonIDCounts,
	}
	canonicalNameUsage.SourceTargetOccurrenceCount.Set(store.DataSourceTypeGBIF, key.TargetID(), 0)

	Ω.Lock()
	defer Ω.Unlock()

	Ω.Usages = append(Ω.Usages, canonicalNameUsage)

	return nil

}

func (Ω *orchestrator) occurrenceCount(id TaxonID) (int, error) {

  	u := fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?limit=1&speciesKey=%d&continent=NORTH_AMERICA&hasCoordinate=true", id)

  	var res struct {
		Count int `json:"count"`
	}

	if err := utils.RequestJSON(u, &res); err != nil {
		return 0, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.OccurrenceCount[id] = res.Count

	return res.Count, nil

}

type MatchResult struct {
	UsageKey       TaxonID    `json:"usageKey"`
	ScientificName string `json:"scientificName"`
	CanonicalName  string `json:"canonicalName"`
	Rank           string `json:"rank"`
	Status         string `json:"status"`
	Confidence     int    `json:"confidence"`
	Note           string `json:"note"`
	MatchType      string `json:"matchType"`
	Synonym    bool   `json:"synonym"`
}



func (Ω *orchestrator) matchSynonyms(id TaxonID) (names, ranks []string, counts name_usage.SourceTargetOccurrenceCount, err error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf( "http://api.gbif.org/v1/species/%d/synonyms?", id))
	if err != nil {
		return nil, nil, nil, err
	}

	ranks = []string{}
	counts = name_usage.SourceTargetOccurrenceCount{}
	names = []string{}

	for _, synonym := range synonymUsages {

		if synonym.Rank == "unranked" {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is unranked, so skipping", synonym.Key))
			continue
		}

		acceptedTaxonomicStatuses := []string{"SYNONYM", "HETEROTYPIC_SYNONYM", "HOMOTYPIC_SYNONYM"}

		if !utils.ContainsString(acceptedTaxonomicStatuses, synonym.TaxonomicStatus) || !synonym.Synonym {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is not a synonym [%s], so skipping", synonym.Key, synonym.TaxonomicStatus))
			continue
		}

		if strings.IndexFunc(synonym.CanonicalName, func(r rune) bool {
			return (r < 'A' || r > 'z') && r != ' ' && r != '-'
		}) != -1 {
			fmt.Println(fmt.Sprintf("Warning: name [%s] contains non letter, so skipping", synonym.CanonicalName))
			continue
		}

		names = utils.AddStringToSet(names, strings.ToLower(synonym.CanonicalName))
		ranks = utils.AddStringToSet(ranks, strings.ToLower(synonym.Rank))
		if counts.Contains(store.DataSourceTypeGBIF, synonym.Key.TargetID()) {
			return nil, nil, nil, errors.Newf("Unexpected: have multiple taxonIDs[%d] within synonyms[%d]", synonym.Key, id)
		}

		counts.Set(store.DataSourceTypeGBIF, synonym.Key.TargetID(), 0)

	}

	return names, ranks, counts, nil
}

type NameUsage struct {
		Key                 TaxonID           `json:"key"`
		NameKey             int           `json:"nameKey"`
		TaxonID             string        `json:"taxonID"`
		Kingdom             string        `json:"kingdom"`
		Phylum              string        `json:"phylum"`
		Order               string        `json:"order"`
		Family              string        `json:"family"`
		Genus               string        `json:"genus"`
		Species             string        `json:"species"`
		KingdomKey          int           `json:"kingdomKey"`
		PhylumKey           int           `json:"phylumKey"`
		ClassKey            int           `json:"classKey"`
		OrderKey            int           `json:"orderKey"`
		FamilyKey           int           `json:"familyKey"`
		GenusKey            int           `json:"genusKey"`
		SpeciesKey          TaxonID           `json:"speciesKey"`
		DatasetKey          string        `json:"datasetKey"`
		ParentKey           TaxonID           `json:"parentKey"`
		Parent              string        `json:"parent"`
		AcceptedKey         TaxonID           `json:"acceptedKey"`
		Accepted            string        `json:"accepted"`
		ScientificName      string        `json:"scientificName"`
		CanonicalName       string        `json:"canonicalName"`
		Authorship          string        `json:"authorship"`
		NameType            string        `json:"nameType"`
		Rank                string        `json:"rank"`
		Origin              string        `json:"origin"`
		TaxonomicStatus     string        `json:"taxonomicStatus"`
		NomenclaturalStatus []interface{} `json:"nomenclaturalStatus"`
		Remarks             string        `json:"remarks"`
		NumDescendants      int           `json:"numDescendants"`
		LastCrawled         string        `json:"lastCrawled"`
		LastInterpreted     string        `json:"lastInterpreted"`
		Issues              []string      `json:"issues"`
		Synonym             bool          `json:"synonym"`
		Class               string        `json:"class"`
		SourceTaxonKey      int           `json:"sourceTaxonKey,omitempty"`
		ConstituentKey      string        `json:"constituentKey,omitempty"`
		BasionymKey         int           `json:"basionymKey,omitempty"`
		Basionym            string        `json:"basionym,omitempty"`
		PublishedIn         string        `json:"publishedIn,omitempty"`
		NubKey              int           `json:"nubKey,omitempty"`
}


func fetchNameUsages(url string) ([]*NameUsage, error) {

	offset := 0
	records := []*NameUsage{}

	for {
		var res struct {
			Offset int `json:"offset"`
			Limit int `json:"limit"`
			EndOfRecords bool `json:"endOfRecords"`
			Results []*NameUsage `json:"results"`
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
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
)

type TaxonID int

type TaxonIDs []TaxonID

func (Ω TaxonIDs) IndexOf(id TaxonID) int {
	for i := range Ω {
		if Ω[i] == id {
			return i
		}
	}
	return -1
}

func (Ω TaxonIDs) AddToSet(id TaxonID) TaxonIDs {
	if Ω.IndexOf(id) == -1 {
		return append(Ω, id)
	}
	return Ω
}

type CanonicalNameUsages []CanonicalNameUsage

func (Ω CanonicalNameUsages) IndexOfNames(qNames ...string) int {
	for i := range Ω {
		if utils.IntersectsStrings(append(Ω[i].Synonyms, Ω[i].CanonicalName), qNames) {
			return i
		}
	}
	return -1
}

func (Ω CanonicalNameUsages) Condense() (CanonicalNameUsages, error) {

ResetLoop:
		for {
			changed := false
			for i := range Ω {
				for k := range Ω {
					if k == i {
						continue
					}
					if Ω[i].ShouldCombine(Ω[k]) {
						changed = true
						var err error
						Ω[i], err = Ω[i].Combine(Ω[k])
						if err != nil {
							return nil, err
						}
						Ω = append(Ω[:k], Ω[k+1:]...)
						continue ResetLoop
					}
				}
			}
			if !changed {
				break
			}
		}

	return Ω, nil
}

func (Ω CanonicalNameUsages) IndexOfID(id TaxonID) int {
	for i := range Ω {
		if _, ok := Ω[i].OccurrenceCount[id]; ok {
			return i
		}
	}
	return -1
}

type CanonicalNameUsage struct {
	CanonicalName string `json:",omitempty"`
	Synonyms []string `json:",omitempty"`
	Ranks []string `json:",omitempty"`
	OccurrenceCount map[TaxonID]int `firestore:",omitempty"`
}

func (a *CanonicalNameUsage) ShouldCombine(b CanonicalNameUsage) bool {
	namesEqual := a.CanonicalName == b.CanonicalName && a.CanonicalName != ""
	if namesEqual {
		fmt.Println("should combine on names", a.CanonicalName, b.CanonicalName)
		return true
	}

	sharesSynonyms := utils.IntersectsStrings(a.Synonyms, b.Synonyms)
	if sharesSynonyms {
		fmt.Println("should combine on synonyms", a.CanonicalName, b.CanonicalName)
		return true
	}

	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)
	if bNameIsSynonym || aNameIsSynonym {
		fmt.Println("should combine on one name is the synonym of the other", a.CanonicalName, b.CanonicalName)
		return true
	}
	sharesTaxonID := false
	for taxonID, _ := range a.OccurrenceCount {
		if _, ok := b.OccurrenceCount[taxonID]; ok {
			sharesTaxonID = true
			break
		}
	}
	if sharesTaxonID {
		fmt.Println("should combine on sharedTaxonID", a.CanonicalName, b.CanonicalName)
		return true
	}

	return false
}

func (a *CanonicalNameUsage) Combine(b CanonicalNameUsage) (CanonicalNameUsage, error) {
	c := CanonicalNameUsage{}

	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)

	if bNameIsSynonym && aNameIsSynonym {
		return c, errors.Newf("What is the real name? %s, %s", a.CanonicalName, b.CanonicalName)
	}

	if bNameIsSynonym {
		c.CanonicalName = a.CanonicalName
	} else {
		//} else if aNameIsSynonym {
		c.CanonicalName = b.CanonicalName
	}

	c.Synonyms = utils.RemoveStringDuplicates(append(a.Synonyms, b.Synonyms...))
	c.Ranks = utils.RemoveStringDuplicates(append(a.Ranks, b.Ranks...))

	c.OccurrenceCount = b.OccurrenceCount
	for taxonID, count := range a.OccurrenceCount {
		if _, ok := c.OccurrenceCount[taxonID]; ok {
			c.OccurrenceCount[taxonID] += count
		} else {
			c.OccurrenceCount[taxonID] = count
		}
	}

	return c, nil
}

func (Ω *CanonicalNameUsage) Valid() bool {
	if !utils.IntersectsStrings([]string{"species", "form", "subspecies", "variety"}, Ω.Ranks) {
		return false
	}
	return true
}

type orchestrator struct {
	Usages CanonicalNameUsages
	OccurrenceCount map[TaxonID]int
	sync.Mutex
	Context context.Context
}

func FetchNamesUsages(cxt context.Context, namesToMatch []string, externalKeyCombinations map[int][]string) (CanonicalNameUsages, error) {

	o := orchestrator{
		Usages: CanonicalNameUsages{},
		OccurrenceCount: map[TaxonID]int{},
		Context: cxt,
	}

	// First let's get all the names together and unique.
	names := utils.RemoveStringDuplicates(namesToMatch)
	for _, rNames := range externalKeyCombinations {
		names = utils.AddStringToSet(names, rNames...)
	}
	names = utils.StringsToLower(names...)

	// Recursively fetch them all
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {
				return o.matchName(name)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	for taxonID, _ := range externalKeyCombinations {
		if i := o.Usages.IndexOfID(TaxonID(taxonID)); i == -1 {
			fmt.Println(fmt.Sprintf("Warning: non-existant taxon id [%d] from external key combination", taxonID))
			if err := o.matchKey(TaxonID(taxonID)); err != nil && err != ErrUnsupported {
				return nil, err
			}
		}
	}

	// Now we must ensure all external key combinations are accounted for.
	for taxonID, names := range externalKeyCombinations {
		i := o.Usages.IndexOfID(TaxonID(taxonID))
		if i == -1 {
			continue
		}
		for _, n := range names {
			if o.Usages[i].CanonicalName != strings.ToLower(n) {
				o.Usages[i].Synonyms = append(o.Usages[i].Synonyms, strings.ToLower(n))
			}
		}
	}

	// Extricate all ids.
	ids := TaxonIDs{}
	for i := range o.Usages {
		for taxonID, _ := range o.Usages[i].OccurrenceCount{
			ids = ids.AddToSet(taxonID)
		}
	}

	// Load cache with occurrence counts
	tmb = tomb.Tomb{}
	tmb.Go(func() error {
		for _, _id := range ids {
			id := _id
			tmb.Go(func() error {
				_, err := o.occurrenceCount(id)
				return err
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	for i := range o.Usages {
		for taxonID, _ := range o.Usages[i].OccurrenceCount{
			o.Usages[i].OccurrenceCount[taxonID] = o.OccurrenceCount[taxonID]
		}
	}

	return o.Usages.Condense()
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

	names, ranks, mapTaxonIDCounts, err := Ω.matchSynonyms(TaxonID(matchResult.UsageKey))
	if err != nil {
		return err
	}

	canonicalNameUsage := CanonicalNameUsage{
		CanonicalName: strings.ToLower(matchResult.CanonicalName),
		Synonyms: names,
		Ranks: utils.AddStringToSet(ranks, strings.ToLower(matchResult.Rank)),
		OccurrenceCount: mapTaxonIDCounts,
	}

	canonicalNameUsage.OccurrenceCount[TaxonID(matchResult.UsageKey)] = 0

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

	canonicalNameUsage := CanonicalNameUsage{
		CanonicalName: strings.ToLower(nameUsage.CanonicalName),
		Synonyms: names,
		Ranks: utils.AddStringToSet(ranks, strings.ToLower(nameUsage.Rank)),
		OccurrenceCount: mapTaxonIDCounts,
	}

	canonicalNameUsage.OccurrenceCount[key] = 0

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



func (Ω *orchestrator) matchSynonyms(id TaxonID) (names, ranks []string, counts map[TaxonID]int, err error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf( "http://api.gbif.org/v1/species/%d/synonyms?", id))
	if err != nil {
		return nil, nil, nil, err
	}

	ranks = []string{}
	counts = map[TaxonID]int{}
	names = []string{}

	for _, synonym := range synonymUsages {

		if synonym.Rank == "unranked" {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is unranked, so skipping", synonym.Key))
			continue
		}

		if synonym.TaxonomicStatus != "SYNONYM" && synonym.TaxonomicStatus != "HETEROTYPIC_SYNONYM" || !synonym.Synonym {
			fmt.Println(fmt.Sprintf("Warning: usage [%d] is not a synonym [%s], so skipping", synonym.Key, synonym.TaxonomicStatus))
			continue
		}

		if strings.IndexFunc(synonym.CanonicalName, func(r rune) bool {
			return (r < 'A' || r > 'z') && r != ' '
		}) != -1 {
			fmt.Println(fmt.Sprintf("Warning: name [%s] contains non letter, so skipping", synonym.CanonicalName))
			continue
		}

		names = utils.AddStringToSet(names, strings.ToLower(synonym.CanonicalName))
		if _, ok := counts[TaxonID(synonym.Key)]; ok {
			return nil, nil, nil, errors.Newf("Unexpected: have multiple taxonIDs[%d] within synonyms[%d]", synonym.Key, id)
		} else {
			counts[TaxonID(synonym.Key)] = 0;
		}
		ranks = utils.AddStringToSet(ranks, strings.ToLower(synonym.Rank))
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
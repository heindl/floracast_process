package gbif

import (
	"context"
	"bitbucket.org/heindl/taxa/store"
	"fmt"
	"net/url"
	"bitbucket.org/heindl/taxa/utils"
	"strconv"
	"gopkg.in/tomb.v2"
	"sync"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

type CanonicalNameUsage struct {
	Name string `firestore:",omitempty" json:",omitempty"`
	SynonymOf string `firestore:",omitempty" json:",omitempty"`
	Rank string `firestore:",omitempty" json:",omitempty"`
	SourceID store.DataSourceID   `firestore:",omitempty" json:",omitempty"`
	TargetID      store.DataSourceTargetID `firestore:",omitempty"`
	Synonyms	[]*CanonicalNameUsage `firestore:",omitempty"`
}

func (Ω *CanonicalNameUsage) Valid() bool {
	if !utils.Contains([]string{"species", "form", "subspecies", "variety"}, Ω.Rank) {
		return false
	}
	return true
}

func MatchNames(cxt context.Context, matches ...string) ([]*CanonicalNameUsage, error) {

	res := []*CanonicalNameUsage{}
	locker := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _matchName := range matches {
			matchName := _matchName
			tmb.Go(func() error {
				usages, err := MatchName(cxt, matchName)
				if err != nil {
					return err
				}
				locker.Lock()
				defer locker.Unlock()

				res = append(res, usages...)

				return nil
			})

		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}



func MatchKey(cxt context.Context, usageKey int) ([]*CanonicalNameUsage, error) {
	// Get the reference for the synonym.
	nameUsage := NameUsage{}
	if err := utils.RequestJSON(fmt.Sprintf("http://api.gbif.org/v1/species/%d?", usageKey), &nameUsage); err != nil {
		return nil, err
	}
	if nameUsage.Synonym {
		if nameUsage.AcceptedKey == usageKey {
			return nil, errors.Newf("i misunderstand how this works: %d - %d", usageKey, nameUsage.AcceptedKey)
		}
		return MatchKey(cxt, nameUsage.AcceptedKey)
	}

	parentUsage := CanonicalNameUsage{
		Name: strings.ToLower(nameUsage.CanonicalName),
		Rank: strings.ToLower(nameUsage.Rank),
		SourceID: store.DataSourceIDGBIF,
		TargetID:store.DataSourceTargetID(strconv.Itoa(nameUsage.Key)),
	}

	if !parentUsage.Valid() {
		return nil, nil
	}

	res, err := MatchSynonyms(cxt, parentUsage.Name, nameUsage.Key)
	if err != nil {
		return nil, err
	}

	return append(res, &parentUsage), nil

}

type MatchResult struct {
	UsageKey       int    `json:"usageKey"`
	ScientificName string `json:"scientificName"`
	CanonicalName  string `json:"canonicalName"`
	Rank           string `json:"rank"`
	Status         string `json:"status"`
	Confidence     int    `json:"confidence"`
	Note           string `json:"note"`
	MatchType      string `json:"matchType"`
	Synonym    bool   `json:"synonym"`
}

func MatchName(cxt context.Context, name string) ([]*CanonicalNameUsage, error) {

	url := fmt.Sprintf("http://api.gbif.org/v1/species/match?name=%s&verbose=true", url.QueryEscape(name))

	matchResult := MatchResult{}
	if err := utils.RequestJSON(url, &matchResult); err != nil {
		return nil, err
	}

	if matchResult.Synonym {
		return MatchKey(cxt, matchResult.UsageKey)
	}

	canonicalNameUsage := CanonicalNameUsage{
		SourceID: store.DataSourceIDGBIF,
		Name: strings.ToLower(matchResult.CanonicalName),
		Rank: strings.ToLower(matchResult.Rank),
		TargetID: store.DataSourceTargetID(strconv.Itoa(matchResult.UsageKey)),
	}

	if !canonicalNameUsage.Valid() {
		return nil, nil
	}

	res, err := MatchSynonyms(cxt, canonicalNameUsage.Name, matchResult.UsageKey)
	if err != nil {
		return nil, err
	}

	return append(res, &canonicalNameUsage), nil
}

func MatchSynonyms(cxt context.Context, synonymOf string, usageKey int) ([]*CanonicalNameUsage, error) {

	synonymUsages, err := fetchNameUsages(fmt.Sprintf( "http://api.gbif.org/v1/species/%d/synonyms?", usageKey))
	if err != nil {
		return nil, err
	}

	res := []*CanonicalNameUsage{}

	for _, synonym := range synonymUsages {

		if synonym.TaxonomicStatus != "SYNONYM" || !synonym.Synonym {
			continue
		}

		cnu := CanonicalNameUsage{
			Name: strings.ToLower(synonym.CanonicalName),
			Rank: strings.ToLower(synonym.Rank),
			TargetID: store.DataSourceTargetID(strconv.Itoa(synonym.Key)),
			SourceID: store.DataSourceIDGBIF,
		}

		if synonymOf != strings.ToLower(synonym.CanonicalName) {
			cnu.SynonymOf = strings.ToLower(synonymOf)
		}

		if cnu.Valid() {
			res = append(res, &cnu)
		}
	}

	return res, nil
}

type NameUsage struct {
		Key                 int           `json:"key"`
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
		SpeciesKey          int           `json:"speciesKey"`
		DatasetKey          string        `json:"datasetKey"`
		ParentKey           int           `json:"parentKey"`
		Parent              string        `json:"parent"`
		AcceptedKey         int           `json:"acceptedKey"`
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
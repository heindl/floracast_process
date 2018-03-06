package gbif

import (
	"bitbucket.org/heindl/process/utils"
	"fmt"
)

// A name usage is the usage of a scientific name according to one particular Checklist including
// the GBIF Taxonomic Backbone which is just called nub in this API. Name usages from other checklists
// with names that also exist in the nub will have a taxonKey that points to the related usage in the backbone.
type species int

func (s species) url() string {
	return fmt.Sprintf("http://api.gbif.org/v1/species/%d", s)
}

type nameUsage struct {
	classification
	AccordingTo         string          `json:"accordingTo"`
	Authorship          string          `json:"authorship"`
	CanonicalName       string          `json:"canonicalName"`
	Confidence          int             `json:"confidence"` // From search api.
	Issues              []interface{}   `json:"issues"`
	Key                 taxonID         `json:"key"`
	LastCrawled         string          `json:"lastCrawled"`
	LastInterpreted     string          `json:"lastInterpreted"`
	MatchType           string          `json:"matchType"` // From search api.
	Modified            string          `json:"modified"`
	NameKey             int             `json:"nameKey"`
	NameType            string          `json:"nameType"`
	NomenclaturalStatus []interface{}   `json:"nomenclaturalStatus"`
	Note                string          `json:"note"` // From search api.
	NubKey              int             `json:"nubKey"`
	NumDescendants      int             `json:"numDescendants"`
	Origin              string          `json:"origin"`
	Parent              string          `json:"parent"`
	ParentKey           taxonID         `json:"parentKey"`
	Rank                rank            `json:"rank"`
	ScientificName      string          `json:"scientificName"`
	SourceTaxonKey      int             `json:"sourceTaxonKey,omitempty"`
	Status              string          `json:"status"` // From search api.
	Synonym             bool            `json:"synonym"`
	TaxonID             string          `json:"taxonID"`
	TaxonomicStatus     taxonomicStatus `json:"taxonomicStatus"`
	UsageKey            taxonID         `json:"usageKey"`    // From search api.
	AcceptedKey         taxonID         `json:"acceptedKey"` // From search api.
	DatasetKey          string          `json:"datasetKey"`
	ConstituentKey      string          `json:"constituentKey,omitempty"`
	BasionymKey         int             `json:"basionymKey,omitempty"`
	Basionym            string          `json:"basionym,omitempty"`
	VernacularName      string          `json:"vernacularName,omitempty"`
	Remarks             string          `json:"remarks,omitempty"`
	PublishedIn         string          `json:"publishedIn,omitempty"`
	Class               string          `json:"class,omitempty"`
	References          string          `json:"references,omitempty"`
}

type classification struct {
	Class      string `json:"class" bson:"class"`
	ClassKey   int    `json:"classKey" bson:"classKey"`
	Family     string `json:"family" bson:"family"`
	FamilyKey  int    `json:"familyKey" bson:"familyKey"`
	Genus      string `json:"genus" bson:"genus"`
	GenusKey   int    `json:"genusKey" bson:"genusKey"`
	Kingdom    string `json:"kingdom" bson:"kingdom"`
	KingdomKey int    `json:"kingdomKey" bson:"kingdomKey"`
	Order      string `json:"order" bson:"order"`
	OrderKey   int    `json:"orderKey" bson:"orderKey"`
	Phylum     string `json:"phylum" bson:"phylum"`
	PhylumKey  int    `json:"phylumKey" bson:"phylumKey"`
	Species    string `json:"species" bson:"species"`
	SpeciesKey int    `json:"speciesKey" bson:"speciesKey"`
}

func (s species) Name() (name nameUsage, err error) {
	if err := utils.RequestJSON(s.url(), &name); err != nil {
		return nameUsage{}, err
	}
	return
}

type parsedNameUsage struct {
	AuthorsParsed           bool   `json:"authorsParsed"`
	BracketAuthorship       string `json:"bracketAuthorship"`
	BracketYear             string `json:"bracketYear"`
	CanonicalName           string `json:"canonicalName"`
	CanonicalNameComplete   string `json:"canonicalNameComplete"`
	CanonicalNameWithMarker string `json:"canonicalNameWithMarker"`
	GenusOrAbove            string `json:"genusOrAbove"`
	Key                     int    `json:"key"`
	ScientificName          string `json:"scientificName"`
	SpecificEpithet         string `json:"specificEpithet"`
	Type                    string `json:"type"`
}

func (s species) ParsedName() (name parsedNameUsage, err error) {
	url := fmt.Sprintf("%s/name", s.url())
	if err := utils.RequestJSON(url, &name); err != nil {
		return parsedNameUsage{}, err
	}
	return
}

func (s species) Parents() (parents []nameUsage, err error) {
	return s.nameUsages("parents")
}

type page struct {
	EndOfRecords bool `json:"endOfRecords"`
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
	Count        int  `json:"count"`
}

func (s species) Children() ([]nameUsage, error) {
	return s.nameUsagePage("children")
}

func (s species) Related() ([]nameUsage, error) {
	return s.nameUsagePage("related")
}

func (s species) Synonyms() ([]nameUsage, error) {
	return s.nameUsagePage("synonyms")
}

func (s species) Combinations() (combinations []nameUsage, err error) {
	return s.nameUsages("combinations")
}

type distribution struct {
	Locality       string `json:"locality"`
	LocationID     string `json:"locationId"`
	Remarks        string `json:"remarks"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
}

func (s species) fetchDistributions() (distributions []distribution, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []distribution `json:"results"`
		}
		url := fmt.Sprintf("%s/distributions?offset=%d&limit=50", s.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		distributions = append(distributions, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

type reference struct {
	Citation       string `json:"citation"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Type           string `json:"type"`
}

func (s species) fetchReferences() (references []reference, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []reference `json:"results"`
		}
		url := fmt.Sprintf("%s/references?offset=%d&limit=50", s.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		references = append(references, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

type vernacularName struct {
	Language       string `json:"language"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	VernacularName string `json:"vernacularName"`
}

func (s species) fetchVernacularNames() (names []vernacularName, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []vernacularName `json:"results"`
		}
		url := fmt.Sprintf("%s/vernacularNames?offset=%d&limit=50", s.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		names = append(names, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

func (s species) nameUsagePage(path string) (names []nameUsage, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []nameUsage `json:"results"`
		}
		url := fmt.Sprintf("%s/%s?offset=%d&limit=50", s.url(), path, offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		names = append(names, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

func (s species) nameUsages(path string) (names []nameUsage, err error) {
	url := fmt.Sprintf("%s/%s", s.url(), path)
	if err := utils.RequestJSON(url, &names); err != nil {
		return nil, err
	}
	return
}

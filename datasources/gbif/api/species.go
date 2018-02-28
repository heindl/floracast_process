package api

import (
	"bitbucket.org/heindl/process/utils"
	"fmt"
)

// A name usage is the usage of a scientific name according to one particular Checklist including
// the GBIF Taxonomic Backbone which is just called nub in this API. Name usages from other checklists
// with names that also exist in the nub will have a taxonKey that points to the related usage in the backbone.
type Species int

func (s Species) url() string {
	return fmt.Sprintf("http://api.gbif.org/v1/species/%d", s)
}

type NameUsage struct {
	*Classification     `json:",inline" bson:",inline"`
	AccordingTo         string          `json:"accordingTo"`
	Authorship          string          `json:"authorship"`
	CanonicalName       string          `json:"canonicalName"`
	Confidence          int             `json:"confidence"` // From search api.
	Issues              []interface{}   `json:"issues"`
	Key                 TaxonID         `json:"key"`
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
	ParentKey           TaxonID         `json:"parentKey"`
	Rank                Rank            `json:"rank"`
	ScientificName      string          `json:"scientificName"`
	SourceTaxonKey      int             `json:"sourceTaxonKey,omitempty"`
	Status              string          `json:"status"` // From search api.
	Synonym             bool            `json:"synonym"`
	TaxonID             string          `json:"taxonID"`
	TaxonomicStatus     TaxonomicStatus `json:"taxonomicStatus"`
	UsageKey            TaxonID         `json:"usageKey"`    // From search api.
	AcceptedKey         TaxonID         `json:"acceptedKey"` // From search api.
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

type Classification struct {
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

func (s Species) Name() (name NameUsage, err error) {
	if err := utils.RequestJSON(s.url(), &name); err != nil {
		return NameUsage{}, err
	}
	return
}

type ParsedNameUsage struct {
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

func (s Species) ParsedName() (name ParsedNameUsage, err error) {
	url := fmt.Sprintf("%s/name", s.url())
	if err := utils.RequestJSON(url, &name); err != nil {
		return ParsedNameUsage{}, err
	}
	return
}

func (s Species) Parents() (parents []NameUsage, err error) {
	return s.nameUsages("parents")
}

type page struct {
	EndOfRecords bool `json:"endOfRecords"`
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
	Count        int  `json:"count"`
}

func (s Species) Children() ([]NameUsage, error) {
	return s.nameUsagePage("children")
}

func (s Species) Related() ([]NameUsage, error) {
	return s.nameUsagePage("related")
}

func (s Species) Synonyms() ([]NameUsage, error) {
	return s.nameUsagePage("synonyms")
}

func (s Species) Combinations() (combinations []NameUsage, err error) {
	return s.nameUsages("combinations")
}

type Description struct {
	Description    string `json:"description"`
	Key            int    `json:"key"`
	Language       string `json:"language"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Type           string `json:"type"`
	License        string `json:"license"`
}

func (s Species) Descriptions() (descriptions []Description, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []Description `json:"results"`
		}
		url := fmt.Sprintf("%s/descriptions?offset=%d&limit=50", s.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		descriptions = append(descriptions, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

type Distribution struct {
	Locality       string `json:"locality"`
	LocationID     string `json:"locationId"`
	Remarks        string `json:"remarks"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
}

func (s Species) Distributions() (distributions []Distribution, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []Distribution `json:"results"`
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

type Media struct {
	Format         string `json:"format"`
	Identifier     string `json:"identifier"`
	License        string `json:"license"`
	Publisher      string `json:"publisher"`
	References     string `json:"references"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	RightsHolder   string `json:"rightsHolder"` // From occurrence search
	Created        string `json:"created"`      // From occurrence search
	Creator        string `json:"creator"`      // From occurrence search
}

func (s Species) Media() (media []Media, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []Media `json:"results"`
		}
		url := fmt.Sprintf("%s/media?offset=%d&limit=50", s.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		media = append(media, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

type Reference struct {
	Citation       string `json:"citation"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Type           string `json:"type"`
}

func (s Species) References() (references []Reference, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []Reference `json:"results"`
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

type VernacularName struct {
	Language       string `json:"language"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	VernacularName string `json:"vernacularName"`
}

func (s Species) VernacularNames() (names []VernacularName, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []VernacularName `json:"results"`
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

func (s Species) nameUsagePage(path string) (names []NameUsage, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []NameUsage `json:"results"`
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

func (s Species) nameUsages(path string) (names []NameUsage, err error) {
	url := fmt.Sprintf("%s/%s", s.url(), path)
	if err := utils.RequestJSON(url, &names); err != nil {
		return nil, err
	}
	return
}

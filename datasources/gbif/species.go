package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strconv"
)

// A name usage is the usage of a scientific name according to one particular Checklist including
// the GBIF Taxonomic Backbone which is just called nub in this API. Name usages from other checklists
// with names that also exist in the nub will have a taxonKey that points to the related usage in the backbone.
type species int

func targetIDToSpecies(id datasources.TargetID) (species, error) {
	i, err := id.ToInt()
	if err != nil {
		return species(0), errors.Wrapf(err, "Could not cast GBIF TargetID [%s] as taxonID", id)
	}
	spcs := species(i)
	if !spcs.Valid() {
		return species(0), errors.Newf("Invalid GBIF taxonID [%s]", id)
	}
	return spcs, nil
}

func (Ω species) Valid() bool {
	return Ω != 0
}

func (Ω species) TargetID() datasources.TargetID {
	return datasources.TargetID(strconv.Itoa(int(Ω)))
}

func (Ω species) url() string {
	return fmt.Sprintf("http://api.gbif.org/v1/species/%d?verbose=true", Ω)
}

type nameUsage struct {
	classification
	AccordingTo         string          `json:"accordingTo"`
	Authorship          string          `json:"authorship"`
	CanonicalName       string          `json:"canonicalName"`
	Confidence          int             `json:"confidence"` // From search api.
	Issues              []interface{}   `json:"issues"`
	Key                 species         `json:"key"`
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
	ParentKey           species         `json:"parentKey"`
	Rank                rank            `json:"rank"`
	ScientificName      string          `json:"scientificName"`
	SourceTaxonKey      int             `json:"sourceTaxonKey,omitempty"`
	Status              string          `json:"status"` // From search api.
	Synonym             bool            `json:"synonym"`
	TaxonID             string          `json:"taxonID"`
	TaxonomicStatus     taxonomicStatus `json:"taxonomicStatus"`
	UsageKey            species         `json:"usageKey"`    // From search api.
	AcceptedKey         species         `json:"acceptedKey"` // From search api.
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

func (Ω species) Name() (nameUsage, error) {
	u := nameUsage{}
	if err := utils.RequestJSON(Ω.url(), &u); err != nil {
		return nameUsage{}, err
	}
	return u, nil
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

func (Ω species) ParsedName() (name parsedNameUsage, err error) {
	url := fmt.Sprintf("%s/name", Ω.url())
	if err := utils.RequestJSON(url, &name); err != nil {
		return parsedNameUsage{}, err
	}
	return
}

func (Ω species) Parents() ([]*nameUsage, error) {
	return Ω.nameUsages("parents")
}

type page struct {
	EndOfRecords bool `json:"endOfRecords"`
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
	Count        int  `json:"count"`
}

func (Ω species) Children() ([]*nameUsage, error) {
	return Ω.nameUsages("children")
}

func (Ω species) Related() ([]*nameUsage, error) {
	return Ω.nameUsages("related")
}

func (Ω species) Synonyms() ([]*nameUsage, error) {
	return Ω.nameUsages("synonyms")
}

func (Ω species) Combinations() ([]*nameUsage, error) {
	return Ω.nameUsages("combinations")
}

type distribution struct {
	Locality       string `json:"locality"`
	LocationID     string `json:"locationId"`
	Remarks        string `json:"remarks"`
	Source         string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
}

func (Ω species) fetchDistributions() (distributions []distribution, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []distribution `json:"results"`
		}
		url := fmt.Sprintf("%s/distributions?offset=%d&limit=50", Ω.url(), offset)
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

func (Ω species) fetchReferences() (references []reference, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []reference `json:"results"`
		}
		url := fmt.Sprintf("%s/references?offset=%d&limit=50", Ω.url(), offset)
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

func (Ω species) fetchVernacularNames() (names []vernacularName, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []vernacularName `json:"results"`
		}
		url := fmt.Sprintf("%s/vernacularNames?offset=%d&limit=50", Ω.url(), offset)
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

func (Ω species) nameUsages(path string) (names []*nameUsage, err error) {

	offset := 0
	records := []*nameUsage{}

	for {
		var res struct {
			page
			Results []*nameUsage `json:"results"`
		}

		if err := utils.RequestJSON(
			fmt.Sprintf("%s/%s?offset=%d&limit=300", Ω.url(), path, offset),
			&res,
		); err != nil {
			return nil, err
		}

		records = append(records, res.Results...)
		if res.EndOfRecords {
			break
		}
		offset += res.Limit
	}

	return records, nil
}

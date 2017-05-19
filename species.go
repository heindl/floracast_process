package species

import (
	"time"
	"github.com/heindl/gbif"
	"fmt"
	"strings"
	"strconv"
	"github.com/saleswise/errors/errors"
	. "github.com/saleswise/malias"
	"bitbucket.org/heindl/utils"
)

type Species struct {
	CanonicalName CanonicalName `bson:"canonicalName,omitempty" json:"canonicalName,omitempty"` // Canonical Name
	ScientificName string `bson:"scientificName,omitempty" json:"scientificName,omitempty"`
	CommonName string `bson:"commonName,omitempty" json:"commonName,omitempty"`
	Sources       map[SourceKey]SourceData       `bson:"sources,omitempty" json:"sources,omitempty"`
	ModifiedAt    *time.Time    `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	CreatedAt    *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	Description   *Media        `bson:"description,omitempty" json:"description,omitempty"`
	Image         *Media        `bson:"image,omitempty" json:"image,omitempty"`
	Classification *gbif.Classification   `bson:"classification,omitempty" json:"classification,omitempty"`
}

type SourceKey string
const sourceKeySeperator = "-+-"

func (this SourceKey) Valid() bool {
	return strings.Contains(string(this), sourceKeySeperator)
}

func NewSourceKey(k IndexKey, t SourceType) SourceKey {
	return SourceKey(fmt.Sprintf("%s%s%v", t, sourceKeySeperator, k))
}

func (this SourceKey) Unmarshal() (IndexKey, SourceType, error) {
	s := strings.Split(string(this), sourceKeySeperator)
	// Right now we only have integer src keys so convert to integer. In future will consider source type before initial conversion.
	i, err := strconv.Atoi(s[1])
	if err != nil {
		return IndexKey(""), SourceType(""), errors.Wrap(err, "could not convert id string to integer").SetState(M{utils.LogkeyStringValue: s[1]})
	}
	return IndexKey(i), SourceType(s[0]), nil
}

type SourceKeys []SourceKey

func (keys SourceKeys) AddToSet(key SourceKey) SourceKeys {
	for _, k := range keys {
		if k == key {
			return keys
		}
	}
	return append(keys, key)
}

type SourceData struct{
	LastFetchedAt *time.Time
}

type SpeciesList []Species

func (list SpeciesList) HasSourceKey(key SourceKey) bool {
	for _, s := range list {
		if s.HasSourceKey(key) {
			return true
		}
	}
	return false
}

func (list SpeciesList) AddToSet(s Species) SpeciesList {
	for _, l := range list {
		if l.CanonicalName == s.CanonicalName {
			return list
		}
	}
	return append(list, s)
}

func (sp *Species) HasSourceKey(key SourceKey) bool {
	if !key.Valid() {
		return false
	}
	if sp == nil {
		return false
	}
	if sp.Sources == nil || len(sp.Sources) == 0 {
		return false
	}
	_, ok := sp.Sources[key]
	return ok
}

// The canonical name is the scientific name of the species, which can cover multiple subspecies.
type CanonicalName string

func (this CanonicalName) Valid() bool {
	return this != ""
}

type Media struct {
	Source string `json:"source" bson:"source"`
	Value  string `json:"value" bson:"value"`
}

type SourceType string

func (s SourceType) Valid() bool {
	return s == SourceTypeGBIF || s == SourceTypeEOL
}

const (
	SourceTypeGBIF = SourceType("gbif")
	SourceTypeEOL  = SourceType("eol")
)

type IndexKey interface{} // Could be a string or an int

package species

import (
	"time"
)

type Species struct {
	CanonicalName CanonicalName `bson:"canonicalName,omitempty" json:"canonicalName,omitempty"` // Canonical Name
	CommonName CanonicalName `bson:"commonName,omitempty" json:"commonName,omitempty"`
	Sources       Sources       `bson:"sources,omitempty" json:"sources,omitempty"`
	ModifiedAt    *time.Time    `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	Description   *Media        `bson:"description,omitempty" json:"description,omitempty"`
	Image         *Media        `bson:"image,omitempty" json:"image,omitempty"`
}

type SpeciesList []Species

func (list SpeciesList) HasSource(src Source) bool {
	for _, l := range list {
		if l.HasSource(src) {
			return true
		}
	}
	return false
}

func (list SpeciesList) AppendUnique(s Species) SpeciesList {
	for _, l := range list {
		if l.CanonicalName == s.CanonicalName {
			return list
		}
	}
	list = append(list, s)
	return list
}

func (sp *Species) HasSource(src Source) bool {
	if src.IsZero() {
		return false
	}
	for _, so := range sp.Sources {
		if so.IndexKey == src.IndexKey && so.Type == src.Type {
			return true
		}
	}
	return false
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
	return string(s) != ""
}

const (
	SourceTypeGBIF = SourceType("gbif")
	SourceTypeEOL  = SourceType("eol")
)

type Sources []Source

type IndexKey int

func (i IndexKey) Valid() bool {
	return int(i) != 0
}

type Source struct {
	Type     SourceType `json:"type" bson:"type"`
	IndexKey IndexKey   `json:"indexKey" bson:"indexKey"`
}

func (this Source) IsZero() bool {
	return this.Type == "" || this.IndexKey == 0
}

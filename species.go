package species

import (
	"time"
	"github.com/saleswise/errors/errors"
	"bitbucket.org/heindl/logkeys"
	"github.com/heindl/gbif"
	. "bitbucket.org/heindl/malias"
)

type Species struct {
	CanonicalName CanonicalName `bson:"canonicalName,omitempty" json:"canonicalName,omitempty"` // Canonical Name
	ScientificName string `bson:"scientificName,omitempty" json:"scientificName,omitempty"`
	CommonName string `bson:"commonName,omitempty" json:"commonName,omitempty"`
	Sources       Sources       `bson:"sources,omitempty" json:"sources,omitempty"`
	ModifiedAt    *time.Time    `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	CreatedAt    *time.Time    `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	Description   *Media        `bson:"description,omitempty" json:"description,omitempty"`
	Image         *Media        `bson:"image,omitempty" json:"image,omitempty"`
	Classification *gbif.Classification   `bson:"classification,omitempty" json:"classification,omitempty"`
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
		if so.IndexKey == src.IndexKey && so.SourceType == src.SourceType {
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
	return s == SourceTypeGBIF || s == SourceTypeEOL
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
	SourceType SourceType `json:"type" bson:"type"`
	IndexKey   IndexKey   `json:"indexKey" bson:"indexKey"`
}

func (this Sources) AddToSet(src SourceType, key IndexKey) (Sources, error) {
	if !src.Valid() {
		return nil, errors.New("invalid source type").SetState(M{logkeys.Source: src})
	}
	for _, s := range this {
		if s.IndexKey == key && s.SourceType == src {
			return this, nil
		}
	}
	return append(this, Source{src, key}), nil
}

func (this Source) IsZero() bool {
	return this.SourceType == "" || this.IndexKey == 0
}
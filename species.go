package species

import (
	"time"
)

type Species struct {
	CanonicalName CanonicalName `bson:"canonicalName,omitempty" json:"canonicalName,omitempty"` // Canonical Name
	Sources       Sources       `bson:"sources,omitempty" json:"sources,omitempty"`
	ModifiedAt    *time.Time    `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	Description   *Media        `bson:"description,omitempty" json:"description,omitempty"`
	Image         *Media        `bson:"image,omitempty" json:"image,omitempty"`
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

const (
	SourceTypeGBIF = SourceType("gbif")
	SourceTypeEOL  = SourceType("eol")
)

type Sources []Source

type IndexKey int

type Source struct {
	Type          SourceType `json:"type" bson:"type"`
	IndexKey IndexKey  `json:"indexKey" bson:"indexKey"`
}

func (this Source) IsZero() bool {
	return this.Type == "" || this.IndexKey == 0
}

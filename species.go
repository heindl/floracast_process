package species

import (
	"bitbucket.org/heindl/cxt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Species struct {
	CanonicalName CanonicalName `bson:"canonicalName,omitempty" json:"canonicalName,omitempty"` // Canonical Name
	Sources       Sources       `bson:"sources,omitempty" json:"sources,omitempty"`
	LastModified  *time.Time    `bson:"lastModified,omitempty" json:"lastModified,omitempty"`
	Description   *Media        `bson:"description,omitempty" json:"description,omitempty"`
	Image         *Media        `bson:"image,omitempty" json:"image,omitempty"`
}

type CanonicalName string
type IndexKey interface{}

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

type Source struct {
	Type     SourceType `json:"type" bson:"type"`
	IndexKey IndexKey   `json:"indexKey" bson:"indexKey"`
}

func (this Source) IsZero() bool {
	return this.Type == "" || this.IndexKey == nil
}

func init() {
	cxt.RegisterCollection(cxt.CollectionInitiator{
		Name: cxt.SpeciesColl,
	})
}

func AllCanonicalNames(c *cxt.Context) (response []CanonicalName, err error) {
	var list []Species
	if err := c.Mongo.Coll(cxt.SpeciesColl).Find(bson.M{}).Select(bson.M{"name": 1}).All(&list); err != nil {
		return nil, errors.Wrap(err, "could not find taxon list")
	}
	for _, s := range list {
		response = append(response, s.CanonicalName)
	}
	return
}

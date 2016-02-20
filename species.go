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
type IndexKey int

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
	return this.Type == "" || this.IndexKey == 0
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

type Fetcher interface{
	Get(CanonicalName) (*Species, error)
}

func NewMockFetcher(list []Species) Fetcher {
	m := make(map[CanonicalName]Species)
	for _, l := range list {
		m[l.CanonicalName] = l
	}
	return Fetcher(&mockfetcher{m})
}

type mockfetcher struct {
	m map[CanonicalName]Species
}

func (this *mockfetcher) Get(n CanonicalName) (*Species, error) {
	if s, ok := this.m[n]; ok {
		return &s, nil
	}
	return nil, nil
}

type fetcher struct {
	*cxt.Context
}

func (this *fetcher) GetSpecies(name CanonicalName) (*Species, error) {
	s := &Species{
		CanonicalName: name,
	}
	if err := this.Mongo.Coll(cxt.SpeciesColl).Find(s).One(&s); err != nil {
		return nil, errors.Wrapf(err, "could not fetch species: %s", name)
	}
	return s, nil
}



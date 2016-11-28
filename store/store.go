package store

import (
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/mgoeco"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/utils"
	"github.com/facebookgo/mgotest"
	"github.com/saleswise/errors/errors"
	"gopkg.in/mgo.v2"
	"time"
	"bitbucket.org/heindl/logkeys"
)

type SpeciesStore interface {
	Read() ([]species.Species, error)
	ReadFromIndexKey(species.IndexKey) (*species.Species, error)
	ReadFromCanonicalNames(...species.CanonicalName) ([]species.Species, error)
	AddSources(name species.CanonicalName, sources ...species.Source) error
	SetDescription(species.CanonicalName, *species.Media) error
	SetImage(species.CanonicalName, *species.Media) error
	Close()
}

const SpeciesColl = mgoeco.CollectionName("species")

func init() {
	mgoeco.RegisterCollection(mgoeco.CollectionInitiator{
		Name: SpeciesColl,
		Indexes: []mgo.Index{
			// Index: SpeciesColl.0
			{
				Background: true,
				Sparse:     true,
				Key:        []string{"canonicalName"},
				Bits:       26,
			},
		},
	})
}

var _ SpeciesStore = &store{}

func NewSpeciesStore() (SpeciesStore, error) {
	m, err := mgoeco.LiveMongo()
	if err != nil {
		return nil, err
	}
	return SpeciesStore(&store{nil, m}), nil
}

type store struct {
	server *mgotest.Server
	mongo  *mgoeco.Mongo
}

func (this *store) Close() {
	this.mongo.Close()
	if this.server != nil {
		this.server.Stop() // Only relevent for tests.
	}
	return
}

func (this *store) Read() (res []species.Species, err error) {
	if err := this.mongo.Coll(SpeciesColl).Find(M{}).All(&res); err != nil {
		return nil, errors.Wrap(err, "could not get all species")
	}
	return
}

func (this *store) ReadFromIndexKey(species.IndexKey) (*species.Species, error) {
	return nil, nil
}

func (this *store) ReadFromCanonicalNames(names ...species.CanonicalName) ([]species.Species, error) {
	q := M{
		"canonicalName": M{"$in": names},
	}
	var list []species.Species
	if err := this.mongo.Coll(SpeciesColl).Find(q).All(&list); err != nil {
		return nil, errors.Wrap(err, "could not find species from canonical names").SetState(M{logkeys.Query: q})
	}
	return list, nil
}

func (this *store) AddSources(name species.CanonicalName, sources ...species.Source) error {
	// Index: SpeciesColl.0
	if _, err := this.mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$addToSet": M{
			"sources": M{"$each": sources},
		},
		"$set": species.Species{
			ModifiedAt: utils.TimePtr(time.Now()),
		},
	}); err != nil {
		errors.Wrap(err, "could not upsert taxon")
	}
	return nil
}

func (this *store) SetDescription(name species.CanonicalName, media *species.Media) error {
	// Index: SpeciesColl.0
	if _, err := this.mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$set": species.Species{
			Description: media,
			ModifiedAt:  utils.TimePtr(time.Now()),
		},
	}); err != nil {
		errors.Wrap(err, "could not add description")
	}
	return nil
}

func (this *store) SetImage(name species.CanonicalName, media *species.Media) error {
	// Index: SpeciesColl.0
	if _, err := this.mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$set": species.Species{
			Image:      media,
			ModifiedAt: utils.TimePtr(time.Now()),
		},
	}); err != nil {
		errors.Wrap(err, "could not add image")
	}
	return nil
}

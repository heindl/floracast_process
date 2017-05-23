package store

import (
	. "github.com/saleswise/malias"
	"bitbucket.org/heindl/provision/mgoeco"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/utils"
	"github.com/facebookgo/mgotest"
	"github.com/saleswise/errors/errors"
	"gopkg.in/mgo.v2"
	"time"
	"github.com/heindl/gbif"
	"fmt"
	"github.com/jonboulle/clockwork"
)

type SpeciesStore interface {
	Read() (species.SpeciesList, error)
	NewIterator() *mgo.Iter
	ReadFromCanonicalNames(...species.CanonicalName) (species.SpeciesList, error)
	ReadFromSourceKeys(...species.SourceKey) (species.SpeciesList, error)
	AddSource(species.CanonicalName, species.SourceKey) error
	SetSourceLastFetched(key species.SourceKey) error
	SetDescription(species.CanonicalName, *species.Media) error
	SetImage(species.CanonicalName, *species.Media) error
	SetClassification(species.CanonicalName, *gbif.Classification) error
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
			// Index: SpeciesColl.0
			{
				Background: true,
				Sparse:     true,
				Key:        []string{"sources"},
				Bits:       26,
			},
		},
	})
}

var _ SpeciesStore = &store{}

func NewTestSpeciesStore(server *mgotest.Server, m *mgoeco.Mongo) SpeciesStore {
	return SpeciesStore(&store{server, m, clockwork.NewFakeClockAt(time.Now())})
}

func NewSpeciesStore() (SpeciesStore, error) {
	m, err := mgoeco.LiveMongo()
	if err != nil {
		return nil, err
	}
	return SpeciesStore(&store{nil, m, clockwork.NewRealClock()}), nil
}

type store struct {
	server *mgotest.Server
	Mongo  *mgoeco.Mongo
	Clock clockwork.Clock
}

func (this *store) Close() {
	this.Mongo.Close()
	if this.server != nil {
		this.server.Stop() // Only relevent for tests.
	}
	return
}

func (this *store) Read() (res species.SpeciesList, err error) {
	if err := this.Mongo.Coll(SpeciesColl).Find(M{}).Sort("canonicalName").All(&res); err != nil {
		return nil, errors.Wrap(err, "could not get all species")
	}
	return
}

func (this *store) NewIterator() *mgo.Iter {
	return this.Mongo.Coll(SpeciesColl).Find(M{}).Iter()
}

func (this *store) ReadFromSourceKeys(keys ...species.SourceKey) (res species.SpeciesList, err error) {
	for _, key := range keys {
		if res.HasSourceKey(key) {
			continue
		}
		q := M{fmt.Sprintf("sources.%s", key): M{"$exists": true}}
		var s species.Species
		if err := this.Mongo.Coll(SpeciesColl).Find(q).One(&s); err != nil {
			return nil, errors.Wrap(err, "could not find species").SetState(M{utils.LogkeyQuery: q})
		}
		res = res.AddToSet(s)
	}
	return res, nil
}

func (this *store) ReadFromCanonicalNames(names ...species.CanonicalName) (species.SpeciesList, error) {
	q := M{
		"canonicalName": M{"$in": names},
	}
	var list []species.Species
	if err := this.Mongo.Coll(SpeciesColl).Find(q).All(&list); err != nil {
		return nil, errors.Wrap(err, "could not find species from canonical names").SetState(M{utils.LogkeyQuery: q})
	}
	return list, nil
}

func (this *store) SetSourceLastFetched(key species.SourceKey) error {
	q := M{
		fmt.Sprintf("sources.%s", key): M{"$exists": true},
	}
	u := M{
		"$set": M{
			fmt.Sprintf("sources.%s.lastFetchedAt", key): this.Clock.Now(),
			"modifiedAt": this.Clock.Now(),
		},
	}
	if err := this.Mongo.Coll(SpeciesColl).Update(q, u); err != nil {
		return errors.Wrap(err, "could not add new source date").SetState(M{
			utils.LogkeyQuery: q,
			utils.LogkeyUpdate: u,
		})
	}
	return nil
}

func (this *store) AddSource(name species.CanonicalName, key species.SourceKey) error {

	// First ensure the species canonicalName if it doesn't exist.
	q := M{
		"canonicalName": name,
	}
	u := M{
		"$setOnInsert": M{
			"createdAt": utils.TimePtr(this.Clock.Now()),
		},
	}
	if _, err := this.Mongo.Coll(SpeciesColl).Upsert(q, u); err != nil {
		return errors.Wrap(err, "could not upsert species").SetState(M{
			utils.LogkeyQuery: q,
			utils.LogkeyUpdate: u,
		})
	}

	q = M{
		"canonicalName": name,
		fmt.Sprintf("sources.%s", key): M{"$exists": false},
	}
	u = M{
		"$set": M{
			fmt.Sprintf("sources.%s", key): M{
				"lastFetchedAt": utils.TimePtr(time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
			"modifiedAt": utils.TimePtr(this.Clock.Now()),
		},
	}
	// Index: SpeciesColl.0
	if err := this.Mongo.Coll(SpeciesColl).Update(q, u); err != nil && err != mgo.ErrNotFound {
		return errors.Wrap(err, "could not upsert taxon").SetState(M{
			utils.LogkeyQuery: q,
			utils.LogkeyUpdate: u,
		})
	}

	return nil
}

func (this *store) SetClassification(name species.CanonicalName, c *gbif.Classification) error {
	// Index: SpeciesColl.0
	if _, err := this.Mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$set": species.Species{
			Classification: c,
		},
	}); err != nil {
		return errors.Wrap(err, "could not add classification")
	}
	return nil
}

func (this *store) SetDescription(name species.CanonicalName, media *species.Media) error {
	// Index: SpeciesColl.0
	if _, err := this.Mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$set": species.Species{
			Description: media,
			ModifiedAt:  utils.TimePtr(time.Now()),
		},
	}); err != nil {
		return errors.Wrap(err, "could not add description")
	}
	return nil
}

func (this *store) SetImage(name species.CanonicalName, media *species.Media) error {
	// Index: SpeciesColl.0
	if _, err := this.Mongo.Coll(SpeciesColl).Upsert(M{"canonicalName": name}, M{
		"$set": species.Species{
			Image:      media,
			ModifiedAt: utils.TimePtr(time.Now()),
		},
	}); err != nil {
		return errors.Wrap(err, "could not add image")
	}
	return nil
}

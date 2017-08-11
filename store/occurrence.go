package store

import (
	"time"
	"cloud.google.com/go/datastore"
	"github.com/saleswise/errors/errors"
	"context"
)

type Occurrences []*Occurrence

func (Ω Occurrences) Find(k *datastore.Key) *Occurrence {
	for _, o := range Ω {
		if o.Key.Kind != k.Kind {
			continue
		}
		if o.Key.ID != k.ID {
			continue
		}
		// The occurrence parent should be a scheme.
		if o.Key.Parent.Name != k.Parent.Name {
			continue
		}
		// The occurrence grandparent should be a taxon.
		if o.Key.Parent.Parent.ID != o.Key.Parent.Parent.ID {
			continue
		}
		return o
	}
	return nil
}

func (Ω Occurrence) Combine(o *Occurrence) *Occurrence {

	if !o.Key.Incomplete() {
		Ω.Key = o.Key
	}
	if o.Location != nil && o.Location.Valid() {
		Ω.Location = o.Location
	}
	if !o.Date.IsZero() {
		Ω.Date = o.Date
	}
	if o.References != "" {
		Ω.References = o.References
	}
	if o.RecordedBy != "" {
		Ω.RecordedBy = o.RecordedBy
	}
	if !o.CreatedAt.IsZero() && o.CreatedAt.Before(Ω.CreatedAt) {
		Ω.CreatedAt = o.CreatedAt
	}
	if !o.ModifiedAt.IsZero() && o.ModifiedAt.After(Ω.ModifiedAt) {
		Ω.ModifiedAt = o.ModifiedAt
	}
	if o.Elevation != 0 {
		Ω.Elevation = o.Elevation
	}
	return &Ω
}

const EntityKindOccurrence = "Occurrence"

type Occurrence struct {
	Key *datastore.Key `datastore:"__key__,omitempty" json:",omitempty" bson:",omitempty"`
	Location *datastore.GeoPoint `datastore:",omitempty" json:"location,omitempty" bson:"location,omitempty"`
	Date      time.Time `datastore:",omitempty,noindex" bson:"date,omitempty" json:"date,omitempty"`
	References string `datastore:",omitempty,noindex" bson:"references,omitempty" json:"references,omitempty"`
	RecordedBy string `datastore:",omitempty,noindex" bson:"recordedBy,omitempty" json:"recordedBy,omitempty"`
	CreatedAt time.Time `datastore:",omitempty,noindex" bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	ModifiedAt time.Time `datastore:",omitempty,noindex" bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	// A globally unique identifier. Although missing in some cases, will be helpful in identifying source of data.
	OccurrenceID string `datastore:",omitempty,noindex" bson:",omitempty" json:",omitempty"`
	Elevation float64 `datastore:",omitempty,noindex" bson:"elevation,omitempty" json:"elevation,omitempty"`
}

func (Ω *Occurrence) Validate() error {
	if Ω == nil {
		return errors.New("nil occurrence")
	}
	if !Ω.Location.Valid() {
		return errors.New("invalid occurrence location")
	}
	if Ω.Key.Name == "" || Ω.Key.Kind != EntityKindOccurrence {
		return errors.New("invalid key")
	}
	return nil
}

func (Ω *store) SetOccurrences(occurrences Occurrences) error {

	keys := make([]*datastore.Key, len(occurrences))
	for i := range occurrences {
		keys[i] = occurrences[i].Key
	}

	if _, err := Ω.DatastoreClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		found := make(Occurrences, len(keys))
		if err := tx.GetMulti(keys, found); err != nil {
			if multierror, ok := err.(datastore.MultiError); ok {
				for _, me := range multierror {
					if me == datastore.ErrNoSuchEntity {
						continue
					} else if me != nil {
						return errors.Wrap(me, "could not get occurrences")
					}
				}
			} else {
				return errors.Wrap(err, "could not get occurrences")
			}
		}

		for i := range found {
			if found[i] == nil {
				found[i] = occurrences.Find(keys[i])
				found[i].CreatedAt = Ω.Clock.Now()
			} else {
				found[i] = found[i].Combine(occurrences.Find(keys[i]))
			}
			found[i].ModifiedAt = Ω.Clock.Now()
		}

		if _, err := tx.PutMulti(keys, found); err != nil {
			return errors.Wrap(err, "could not update species data source")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil

}

func (Ω *store) GetOccurrenceIterator(taxonKey *datastore.Key) *datastore.Iterator {
	q := datastore.NewQuery(EntityKindOccurrence)
	if taxonKey != nil {
		q = q.Ancestor(taxonKey)
	}
	q = q.Order("__key__")
	return Ω.DatastoreClient.Run(context.Background(), q)
}

func (Ω *store) GetOccurrences(taxonKey *datastore.Key) (res Occurrences, err error) {
	q := datastore.NewQuery(EntityKindOccurrence)
	if taxonKey != nil {
		q = q.Ancestor(taxonKey)
	}
	q = q.Order("__key__")
	if _, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res); err != nil {
		return nil, errors.Wrap(err, "could not get occurrences")
	}
	return res, nil
}
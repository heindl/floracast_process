package store

import (
	"github.com/saleswise/errors/errors"
	"context"
	"cloud.google.com/go/datastore"
	"time"
	"strings"
	"fmt"
	"bitbucket.org/heindl/utils"
	. "github.com/saleswise/malias"
)

type SchemeSourceID string
const (
	SchemeSourceIDGBIF = SchemeSourceID("27")
	SchemeSourceIDINaturalist = SchemeSourceID("F1")
)

var SchemeSourceIDMap = map[SchemeSourceID]string{
	// Floracast
	SchemeSourceIDINaturalist: "iNaturalist",
	// INaturalist
	SchemeSourceID("1"): "IUCN Red List of Threatened Species. Version 2012.1",
	SchemeSourceID("3"): "Amphibian Species of the World 5.6",
	SchemeSourceID("2"): "Amphibiaweb. 2012",
	SchemeSourceID("5"): "Amphibian Species of the World 5.5",
	SchemeSourceID("17"): "New England Wild Flower Society's Flora Novae Angliae",
	SchemeSourceID("11"): "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	SchemeSourceID("12"): "Calflora",
	SchemeSourceID("13"): "Odonata Central",
	SchemeSourceID("14"): "IUCN Red List of Threatened Species. Version 2012.2",
	SchemeSourceID("10"): "eBird/Clements Checklist 6.7",
	SchemeSourceID("15"): "CONABIO",
	SchemeSourceID("6"): "The Reptile Database",
	SchemeSourceID("16"): "Afribats",
	SchemeSourceID("18"): "Norma 059, 2010",
	SchemeSourceID("4"): "Draft IUCN/SSC, 2013.1",
	SchemeSourceID("19"): "Draft IUCN/SSC Amphibian Specialist Group, 2011",
	SchemeSourceID("20"): "eBird/Clements Checklist 6.8",
	SchemeSourceID("21"): "IUCN Red List of Threatened Species. Version 2013.2",
	SchemeSourceID("22"): "eBird/Clements Checklist 6.9",
	SchemeSourceID("23"): "NatureWatch NZ",
	SchemeSourceID("24"): "The world spider catalog, version 15.5",
	SchemeSourceID("25"): "Carabidae of the World",
	SchemeSourceID("26"): "IUCN Red List of Threatened Species. Version 2014.3",
	SchemeSourceIDGBIF: "GBIF",
	SchemeSourceID("28"): "NPSpecies",
	SchemeSourceID("29"): "Esslinger&#39;s North American Lichens",
	SchemeSourceID("30"): "Amphibian Species of the World 6.0",
	SchemeSourceID("31"): "Esslinger&#39;s North American Lichens, Version 21",
}

type SchemeTargetID string

type Scheme struct{
	Key           *datastore.Key `datastore:"__key__"`
	//SourceID      SchemeSourceID `datastore:",omitempty,noindex" json:"sourceID,omitempty" bson:"sourceID,omitempty"`
	//TargetID      SchemeTargetID `datastore:",omitempty,noindex" json:"targetID,omitempty" bson:"targetID,omitempty"`
	CreatedAt     time.Time `datastore:",omitempty,noindex" json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt    time.Time `datastore:",omitempty,noindex" json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	LastFetchedAt time.Time `datastore:",omitempty" json:"lastFetchedAt,omitempty" bson:"lastFetchedAt,omitempty"`
}

type SchemeKeyName string
const schemeKeyNameSeperator  = "|||"

func NewSchemeKeyName(id SchemeSourceID, targetID SchemeTargetID) SchemeKeyName {
	return SchemeKeyName(fmt.Sprintf("%s%s%s", id, schemeKeyNameSeperator, targetID))
}

func (Ω SchemeKeyName) Validate() error {
	if !strings.Contains(string(Ω), schemeKeyNameSeperator) {
		return errors.New("invalid scheme key name").SetState(M{utils.LogkeyScheme: Ω})
	}
	return nil
}

func (Ω *Scheme) Validate() error {

	if err := SchemeKeyName(Ω.Key.Name).Validate(); err != nil {
		return err
	}

	if Ω.Key.Kind != EntityKindMetaScheme && Ω.Key.Kind != EntityKindOccurrenceScheme {
		return errors.New("wrong key kind for scheme")
	}

	if !ValidTaxonKey(Ω.Key.Parent) {
		return errors.New("invalid scheme parent taxon")
	}

	return nil
}

func (Ω Scheme) Combine(s *Scheme) *Scheme {
	if s.CreatedAt.Before(Ω.CreatedAt) {
		Ω.CreatedAt = s.CreatedAt
	}
	if Ω.LastFetchedAt.Before(s.LastFetchedAt) {
		Ω.LastFetchedAt = s.LastFetchedAt
	}
	return &Ω
}

func newScheme(entityKind string, origin SchemeSourceID, target SchemeTargetID, txn *datastore.Key) *Scheme {
	return &Scheme{
			Key: datastore.NameKey(entityKind, string(NewSchemeKeyName(origin, target)), txn),
			CreatedAt: time.Now(),
			ModifiedAt: time.Now(),
	}
}

const EntityKindMetaScheme = "MetaScheme"
func NewMetaScheme(origin SchemeSourceID, id SchemeTargetID, txn *datastore.Key) *Scheme {
	return newScheme(EntityKindMetaScheme, origin, id, txn)
}
const EntityKindOccurrenceScheme = "OccurrenceScheme"
func NewOccurrenceScheme(origin SchemeSourceID, id SchemeTargetID, parent *datastore.Key) *Scheme {
	return newScheme(EntityKindOccurrenceScheme, origin, id, parent)
}

//func (Ω *store) getScheme(kind string, origin SchemeSourceID, id SchemeTargetID, txn *datastore.Key) (*Scheme, error) {
//
//	q := datastore.NewQuery(kind).
//		Filter("SourceID =", string(origin)).
//		Filter("TargetID =", string(id)).
//		Ancestor(txn)
//
//	res := []*Scheme{}
//	keys, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not get datastore scheme")
//	}
//
//	if len(keys) > 1 {
//		return nil, errors.New("unexpectedly found more than one meta scheme").SetState(M{
//			utils.LogkeyIdentifier: id,
//			utils.LogkeyStringValue: origin,
//			utils.LogkeyDatastoreKey: txn,
//		})
//	}
//
//	if len(keys) == 0 {
//		return nil, nil
//	}
//
//	return res[0], nil
//
//}

func (Ω *store) GetOccurrenceSchema(taxonKey *datastore.Key) (Schema, error) {
	q := datastore.NewQuery(EntityKindOccurrenceScheme)
	if taxonKey != nil {
		q = q.Ancestor(taxonKey)
	}
	q = q.Order("__key__")
	res := Schema{}
	if _, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res); err != nil {
		return nil, errors.Wrap(err, "could not fetch occurrence schema")
	}
	return res, nil
}

func (Ω *store) SetSchema(schema Schema) error {

	keys := []*datastore.Key{}
	for _, s := range schema {
		keys = append(keys, s.Key)
	}

	if _, err := Ω.DatastoreClient.PutMulti(context.Background(), keys, schema); err != nil {
		return errors.Wrap(err, "could not batch scheme puts to datastore")
	}
	return nil
}

//func (Ω *store) CreateScheme(s *Scheme) (*datastore.Key, error) {
//
//	if !ValidTaxonKey(s.Key.Parent) {
//		return nil, errors.New("invalid scheme parent taxon key").SetState(M{utils.LogkeyDatastoreKey: s.Key.Parent})
//	}
//	//if s.SourceID == "" {
//	//	return nil, errors.New("invalid scheme id origin")
//	//}
//	//if s.TargetID == "" {
//	//	return nil, errors.New("invalid scheme id")
//	//}
//
//	// No opt if we find the scheme. Can have an additional update function below if
//	// we need to update additional fields.
//	if sch, err := Ω.getScheme(s.Key.Kind, s.SourceID, s.TargetID, s.Key.Parent); err != nil {
//		return nil, err
//	} else if sch != nil {
//		return sch.Key, nil
//	}
//
//	k, err := Ω.DatastoreClient.Put(context.Background(), s.Key, s)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not set meta scheme")
//	}
//
//	return k, nil
//}

func (Ω *store) UpdateSchemaLastFetched(schema Schema) error {

	keys := []*datastore.Key{}
	for _, s := range schema {
		keys = append(keys, s.Key)
	}

	if _, err := Ω.DatastoreClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		found := Schema{}
		if err := tx.GetMulti(keys, &found); err != nil {
			return errors.Wrap(err, "could not get species data source")
		}

		if len(found) != len(keys) {
			return errors.New("number of records found not equal to records requested")
		}

		nkeys := []*datastore.Key{}
		for i := range found {
			nkeys = append(nkeys, found[i].Key) // TODO: This step may be unnecessary as the records may be returned in the same order as requested. No time to check now.
			found[i] = found[i].Combine(schema.Find(found[i].Key))
			found[i].ModifiedAt = Ω.Clock.Now()
		}

		if _, err := tx.PutMulti(nkeys, found); err != nil {
			return errors.Wrap(err, "could not update species data source")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type Schema []*Scheme

func (Ω Schema) Find(k *datastore.Key) *Scheme {
	for _, s := range Ω {
		if s.Key.Name == k.Name && s.Key.Parent.ID == k.Parent.ID {
			return s
		}
	}
	return nil
}

func (Ω Schema) AddToSet(s *Scheme) (Schema, error) {

	if err := s.Validate(); err != nil {
		return nil, err
	}

	for i := range Ω {
		if Ω[i].Key.Kind != s.Key.Kind {
			continue
		}
		if Ω[i].Key.Name != s.Key.Name {
			continue
		}
		if Ω[i].Key.Parent.ID != s.Key.Parent.ID {
			continue
		}
		Ω[i] = Ω[i].Combine(s)
		return Ω, nil
	}

	return append(Ω, s), nil

}

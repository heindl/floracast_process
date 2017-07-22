package store

import (
	"time"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"github.com/saleswise/errors/errors"
)

const EntityKindTaxon = "Taxon"

func ValidTaxonKey(k *datastore.Key) bool {
	return k != nil && k.Kind == EntityKindTaxon && k.ID != 0
}

// TaxonID is the iNaturalist taxon_id. In order to move quickly we'll be using their taxonomy
// structure which is kept up to date by a network of volunteers using other sources.
// It is an int64 becasu that is required for a
type TaxonID int64

func (Ω TaxonID) Valid() bool {
	return Ω != 0
}

func NewTaxonKey(id int64) *datastore.Key {
	if !TaxonID(id).Valid() {
		return nil
	}
	return datastore.IDKey(EntityKindTaxon, id, nil)
}

type TaxonKeys []*datastore.Key

func (Ω TaxonKeys) Find(id int64, kind string) *datastore.Key {
	for _, k := range Ω {
		if k.ID == id && k.Kind == kind {
			return k
		}
	}
	return nil
}

type RankLevel int
const (
	// Originating from INaturalist:
	RankLevelKingdom = RankLevel(70)
	RankLevelPhylum = RankLevel(60)
	RankLevelSubPhylum = RankLevel(57)
	RankLevelClass = RankLevel(50)
	RankLevelSubClass = RankLevel(47)
	RankLevelOrder = RankLevel(40)
	RankLevelSuperFamily = RankLevel(33)
	RankLevelFamily = RankLevel(30)
	RankLevelSubFamily = RankLevel(27)
	RankLevelTribe = RankLevel(25)
	RankLevelSubTribe = RankLevel(24)
	RankLevelGenus = RankLevel(20)
	RankLevelSpecies = RankLevel(10)
	RankLevelSubSpecies = RankLevel(5)
)

type Taxon struct {
	Key             *datastore.Key `datastore:"__key__"`
	CanonicalName CanonicalName `datastore:",omitempty" bson:"canonicalName,omitempty" json:"canonicalName,omitempty"`
	Rank	TaxonRank `datastore:",omitempty,noindex" bson:"rank,omitempty" json:"rank,omitempty"`
	RankLevel RankLevel `datastore:"RankLevel,omitempty" bson:"rankLevel,omitempty" json:"rankLevel,omitempty"`
	CommonName string `datastore:",omitempty,noindex" bson:"commonName,omitempty" json:"commonName,omitempty"`
	ModifiedAt      time.Time    `datastore:",omitempty,noindex" bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	CreatedAt       time.Time    `datastore:",omitempty,noindex" bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	States []State `datastore:",omitempty" bson:"states,omitempty" json:"states,omitempty"`
	WikipediaSummary string `datastore:",omitempty,noindex" bson:"wikipediaSummary,omitempty" json:"wikipediaSummary,omitempty"`
}

func (Ω Taxon) Combine(s *Taxon) *Taxon {

	if s.Key.Parent != nil && TaxonID(s.Key.Parent.ID).Valid() {
		Ω.Key.Parent = s.Key.Parent
	}

	if s.CreatedAt.Before(Ω.CreatedAt) {
		Ω.CreatedAt = s.CreatedAt
	}
	if s.CanonicalName.Valid() {
		Ω.CanonicalName = s.CanonicalName
	}
	if s.Rank.Valid() {
		Ω.Rank = s.Rank
	}
	if s.RankLevel != 0 {
		Ω.RankLevel = s.RankLevel
	}
	if s.CommonName == "" {
		Ω.CommonName = s.CommonName
	}
	if len(s.States) > 0 {
		Ω.States = s.States
	}

	if s.WikipediaSummary != "" {
		Ω.WikipediaSummary = s.WikipediaSummary
	}

	return &Ω
}

type State struct {
	EstablishmentMeans string `datastore:",omitempty,noindex" bson:"establishmentMeans,omitempty" json:"establishmentMeans,omitempty"`
	Name               string `datastore:",omitempty" bson:"name,omitempty" json:"name,omitempty"`
}

type TaxonRank string

func (Ω TaxonRank) Valid() bool {
	return Ω != ""
}

type Taxa []*Taxon

func (Ω Taxa) AddToSet(s *Taxon) (Taxa, error) {

	if !TaxonID(s.Key.ID).Valid() {
		return nil, errors.New("invalid taxon id")
	}

	for i := range Ω {
		if Ω[i].Key.ID != s.Key.ID {
			continue
		}
		Ω[i] = Ω[i].Combine(s)
		return Ω, nil
	}
	return append(Ω, s), nil
}

func (Ω Taxa) Index(id TaxonID) int {
	for i := range Ω {
		if TaxonID(Ω[i].Key.ID) == id {
			return i
		}
	}
	return -1
}

// The canonical name is the scientific name of the species, which can cover multiple subspecies.
type CanonicalName string

type CanonicalNames []CanonicalName

func (list CanonicalNames) AddToSet(s CanonicalName) CanonicalNames {
	for _, l := range list {
		if l == s {
			return list
		}
	}
	return append(list, s)
}

func (Ω CanonicalName) Valid() bool {
	return Ω != ""
}

func (Ω *store) SetTaxa(txa Taxa) error {
	keys := make([]*datastore.Key, len(txa))
	for i := range txa {
		keys[i] = txa[i].Key
	}
	if _, err := Ω.DatastoreClient.PutMulti(context.Background(), keys, txa); err != nil {
		return errors.Wrap(err, "could not set taxa")
	}
	return nil
}

func (Ω *store) ReadTaxa() (res Taxa, err error) {
	// Filter only to species and subspecies.
	//q := datastore.NewQuery(EntityKindTaxon).Filter("RankLevel =", 5).Filter("RankLevel =", 10).Order("__key__")
	q := datastore.NewQuery(EntityKindTaxon)
	if _, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res); err != nil {
		return nil, errors.Wrap(err, "could not fetch taxa from datastore")
	}
	return res, nil
}

func (Ω *store) NewIterator() *datastore.Iterator {
	// Filter only to species and subspecies.
	return Ω.DatastoreClient.Run(context.Background(), datastore.NewQuery(EntityKindTaxon).Filter("RankLevel =", 5).Filter("RankLevel =", 10))
}

func (Ω *store) ReadTaxaFromCanonicalNames(names ...CanonicalName) (Taxa, error) {
	q := datastore.NewQuery(EntityKindTaxon)
	for _, name := range names {
		q = q.Filter("CanonicalName =", string(name))
	}
	res := Taxa{}
	if _, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res); err != nil {
		return nil, errors.Wrap(err, "could not get taxa from canonical names")
	}
	return res, nil
}

func (Ω *store) GetTaxon(k *datastore.Key) (*Taxon, error) {
	if k.Kind != string(EntityKindTaxon) {
		return nil, errors.New("invalid entity kind")
	}
	res := Taxon{}
	if err := Ω.DatastoreClient.Get(context.Background(), k , &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (Ω *store) ReadSpecies() (Taxa, error) {
	// Filter only to species and subspecies.
	q := datastore.NewQuery(EntityKindTaxon).Filter("RankLevel <=", int(RankLevelSpecies))
	var res Taxa
	if _, err := Ω.DatastoreClient.GetAll(context.Background(), q, &res); err != nil {
		return nil, errors.Wrap(err, "could not fetch species from datastore")
	}
	return res, nil
}

//func (Ω *store) CreateTaxon(txn *Taxon, shouldHaveParent bool) (*datastore.Key, error) {
//	if txn == nil || !ValidTaxonKey(txn.Key) {
//		return nil, errors.New("invalid taxon").SetState(M{utils.LogkeyTaxon: txn})
//	}
//	if shouldHaveParent {
//		if !ValidTaxonKey(txn.Key.Parent) {
//			return nil, errors.New("invalid parent taxon key").SetState(M{utils.LogkeyTaxon: txn})
//		}
//	}
//
//	if !txn.CanonicalName.Valid() {
//		return nil, errors.New("taxon name invalid").SetState(M{utils.LogkeyTaxon: txn})
//	}
//
//	if !txn.Rank.Valid() {
//		return nil, errors.New("taxon rank invalid").SetState(M{utils.LogkeyTaxon: txn})
//	}
//
//	// No-op if it the taxon already exists. We're just creating here.
//	var t Taxon
//	if err := Ω.DatastoreClient.Get(context.Background(), txn.Key, &t); err != nil && err != datastore.ErrNoSuchEntity {
//		return nil, errors.Wrap(err, "could not get taxon").SetState(M{utils.LogkeyDatastoreKey: txn.Key.Kind})
//	} else if err == nil {
//		// We have the entity, so exit early.
//		return t.Key, nil
//	}
//
//	k, err := Ω.DatastoreClient.Put(context.Background(), txn.Key, txn)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not set taxon").SetState(M{utils.LogkeyDatastoreKey: txn.Key.Kind})
//	}
//	return k, nil
//}

// AddTaxonAdminAreas adds geographical areas that contain Ω taxon.
//func (Ω *store) AddTaxonAdminAreas(k *TaxonKey, states ...string) error {
//	if !k.Valid() {
//		return errors.New("invalid taxon key")
//	}
//	_, err := Ω.DatastoreClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
//		src := Taxon{}
//		if err := tx.Get(k, src); err != nil {
//			return errors.Wrap(err, "could not get species data source")
//		}
//		src.States = append(src.States, states...)
//		src.ModifiedAt = Ω.Clock.Now()
//		if src.CreatedAt.IsZero() {
//			src.CreatedAt = Ω.Clock.Now()
//		}
//		if _, err := tx.Put(src.Key, src); err != nil {
//			return errors.Wrap(err, "could not update species data source")
//		}
//		return nil
//	})
//	return err
//}
//
//// AddTaxonCommonName adds a common name to taxon
//func (Ω *store) AddTaxonCommonName(k *TaxonKey, name string) error {
//	if !k.Valid() {
//		return errors.New("invalid taxon key")
//	}
//	_, err := Ω.DatastoreClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
//		src := Taxon{}
//		if err := tx.Get(k, src); err != nil {
//			return errors.Wrap(err, "could not get species data source")
//		}
//		src.CommonName = name
//		src.ModifiedAt = Ω.Clock.Now()
//		if src.CreatedAt.IsZero() {
//			src.CreatedAt = Ω.Clock.Now()
//		}
//		if _, err := tx.Put(src.Key, src); err != nil {
//			return errors.Wrap(err, "could not update species data source")
//		}
//		return nil
//	})
//	return err
//}
//
//// AddTaxonCommonName adds a common name to taxon
//func (Ω *store) AddTaxonWikipediaSummary(k *TaxonKey, summary string) error {
//	if !k.Valid() {
//		return errors.New("invalid taxon key")
//	}
//	_, err := Ω.DatastoreClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
//		src := Taxon{}
//		if err := tx.Get(k, src); err != nil {
//			return errors.Wrap(err, "could not get species data source")
//		}
//		src.WikipediaSummary = summary
//		src.ModifiedAt = Ω.Clock.Now()
//		if src.CreatedAt.IsZero() {
//			src.CreatedAt = Ω.Clock.Now()
//		}
//		if _, err := tx.Put(src.Key, src); err != nil {
//			return errors.Wrap(err, "could not update species data source")
//		}
//		return nil
//	})
//	return err
//}


package store

import (
	"time"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"github.com/saleswise/errors/errors"
	"strings"
	"github.com/fatih/structs"
	"gopkg.in/mgo.v2/txn"
)

//const EntityKindTaxon = "Taxon"
//
//func ValidTaxonKey(k *datastore.Key) bool {
//	return k != nil && TaxonRank(k.Kind).Valid() && k.ID != 0
//}

// ID is the iNaturalist taxon_id. In order to move quickly we'll be using their taxonomy
// structure which is kept up to date by a network of volunteers using other sources.
type TaxonID string

func (Ω TaxonID) Valid() bool {
	return Ω != ""
}

//func NewTaxonKey(id int64, rank TaxonRank) *datastore.Key {
//	if !ID(id).Valid() {
//		return nil
//	}
//	if !rank.Valid() {
//		return nil
//	}
//	return datastore.IDKey(string(rank), id, nil)
//}
//
//type TaxonKeys []*datastore.Key
//
//func (Ω TaxonKeys) Find(id int64, kind string) *datastore.Key {
//	for _, k := range Ω {
//		if k.ID == id && k.Kind == kind {
//			return k
//		}
//	}
//	return nil
//}

type TaxonRank string
const (
	// Originating from INaturalist:
	RankKingdom = TaxonRank("Kingdom")
	RankPhylum = TaxonRank("Phylum")
	RankSubPhylum = TaxonRank("SubPhylum")
	RankClass = TaxonRank("Class")
	RankSubClass = TaxonRank("SubClass")
	RankOrder = TaxonRank("Order")
	RankSuperFamily = TaxonRank("SuperFamily")
	RankFamily = TaxonRank("Family")
	RankSubFamily = TaxonRank("SubFamily")
	RankTribe = TaxonRank("Tribe")
	RankSubTribe = TaxonRank("SubTribe")
	RankGenus = TaxonRank("Genus")
	RankSpecies = TaxonRank("Species")
	RankSubSpecies = TaxonRank("SubSpecies")
	RankForm = TaxonRank("Form")
)

var TaxonRankMap = map[string]TaxonRank{
	"kingdom": RankKingdom,
	"phylum": RankPhylum,
	"subphylum": RankSubPhylum,
	"class": RankClass,
	"subclass": RankSubClass,
	"order": RankOrder,
	"superfamily": RankSuperFamily,
	"family": RankFamily,
	"subfamily": RankSubFamily,
	"tribe": RankTribe,
	"subtribe": RankSubTribe,
	"genus": RankGenus,
	"species": RankSpecies,
	"subspecies": RankSubSpecies,
	"form": RankForm,
}

func (Ω TaxonRank) Valid() bool {
	if _, ok := TaxonRankMap[strings.ToLower(string(Ω))]; !ok {
		return false
	}
	return true
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
	CanonicalName    CanonicalName `firestore:",omitempty"`
	ID               TaxonID       `firestore:",omitempty"`
	ParentID         TaxonID       `firestore:",omitempty"`
	PhotoURL string `firestore:",omitempty"`
	Rank             TaxonRank     `firestore:",omitempty"`
	RankLevel        RankLevel     `firestore:",omitempty"`
	CommonName       string        `firestore:",omitempty"`
	ModifiedAt       time.Time     `firestore:",omitempty"`
	CreatedAt        time.Time     `firestore:",omitempty"`
	States           []State       `firestore:",omitempty"`
	WikipediaSummary string        `firestore:",omitempty"`
}

//func (Ω Taxon) Combine(s *Taxon) *Taxon {
//
//	if s.Key.Parent != nil && ID(s.Key.Parent.ID).Valid() {
//		Ω.Key.Parent = s.Key.Parent
//	}
//
//	if s.CreatedAt.Before(Ω.CreatedAt) {
//		Ω.CreatedAt = s.CreatedAt
//	}
//	if s.CanonicalName.Valid() {
//		Ω.CanonicalName = s.CanonicalName
//	}
//
//	if ID(s.ID).Valid() {
//		Ω.ID = s.ID
//	}
//
//	if s.Rank.Valid() {
//		Ω.Rank = s.Rank
//	}
//	if s.RankLevel != 0 {
//		Ω.RankLevel = s.RankLevel
//	}
//	if s.CommonName == "" {
//		Ω.CommonName = s.CommonName
//	}
//	if len(s.States) > 0 {
//		Ω.States = s.States
//	}
//
//	if s.WikipediaSummary != "" {
//		Ω.WikipediaSummary = s.WikipediaSummary
//	}
//
//	return &Ω
//}

type State struct {
	EstablishmentMeans string `datastore:",omitempty,noindex" bson:"establishmentMeans,omitempty" json:"establishmentMeans,omitempty"`
	Name               string `datastore:",omitempty" bson:"name,omitempty" json:"name,omitempty"`
}

type Taxa []Taxon

func (Ω Taxa) RemoveDuplicates() (response Taxa) {
	for _, t := range Ω {
		if response.Find(t.ID) == nil {
			response = append(response, t)
		}
	}
	return
}

func (Ω Taxa) Find(k TaxonID) *Taxon {
	for _, t := range Ω {
		if t.ID == k {
			return &t
		}
	}
	return nil
}

//func (Ω Taxa) AddToSet(s *Taxon) (Taxa, error) {
//
//	if !ID(s.ID).Valid() {
//		return nil, errors.New("invalid taxon id")
//	}
//
//	for i := range Ω {
//		if Ω[i].Key.Kind != s.Key.Kind {
//			continue
//		}
//		if Ω[i].Key.ID != s.Key.ID {
//			continue
//		}
//		Ω[i] = Ω[i].Combine(s)
//		return Ω, nil
//	}
//	return append(Ω, s), nil
//}

func (Ω Taxa) Index(id TaxonID) int {
	for i := range Ω {
		if TaxonID(Ω[i].ID) == id {
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

func (Ω *store) UpsertTaxon(cxt context.Context, txn Taxon) error {

	// Validate
	if !txn.ID.Valid() {
		return errors.New("invalid taxa")
	}

	ref := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(txn.ID))

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		return tx.UpdateMap(ref, structs.Map(txn))
	}); err != nil {
		return errors.Wrap(err, "could not update occurrence")
	}

	return nil
}

func (Ω *store) SetTaxonPhoto(cxt context.Context, taxonID TaxonID, photoURL string) error {
	ref := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(taxonID))

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		return tx.UpdateMap(ref, map[string]interface{}{
			"PhotoURL": photoURL,
		})
	}); err != nil {
		return errors.Wrap(err, "could not update taxon photo")
	}
}

//func (Ω *store) ReadTaxa() (res Taxa, err error) {
//	// Filter only to species and subspecies.
//	//q := datastore.NewQuery(EntityKindTaxon).Filter("RankLevel =", 5).Filter("RankLevel =", 10).Order("__key__")
//	q := datastore.NewQuery(EntityKindTaxon)
//	if _, err := Ω.FirebaseClient.GetAll(context.Background(), q, &res); err != nil {
//		return nil, errors.Wrap(err, "could not fetch taxa from datastore")
//	}
//	return res, nil
//}

func (Ω *store) ReadTaxaFromCanonicalNames(cxt context.Context, rank TaxonRank, names ...CanonicalName) (res Taxa, err error) {

	q := Ω.FirestoreClient.Collection(CollectionTypeTaxa).
		Where("Rank", "==", rank)

	for _, name := range names {
		q = q.Where("CanonicalName", "==", name)
	}

	docs, err := q.Documents(cxt).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get taxa")
	}

	for _, doc := range docs {
		t := Taxon{}
		if err := doc.DataTo(&t); err != nil {
			return nil, errors.Wrap(err, "could not type cast taxon")
		}
		res = append(res, t)
	}
	return

}

func (Ω *store) GetTaxon(cxt context.Context, id TaxonID) (*Taxon, error) {

	if !id.Valid() {
		return nil, errors.New("invalid taxon id")
	}

	snap, err := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(id)).Get(cxt)
	if err != nil {
		return nil, errors.Wrap(err, "could not get taxon")
	}

	t := Taxon{}
	if err := snap.DataTo(&t); err != nil {
		return nil, errors.Wrap(err, "could not type cast taxon")
	}
	return &t, nil
}

func (Ω *store) ReadTaxa(cxt context.Context) (res Taxa, err error) {

	snapshots, err := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Documents(cxt).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get taxa")
	}

	for _, doc := range snapshots {
		t := Taxon{}
		if err := doc.DataTo(&t); err != nil {
			return nil, errors.Wrap(err, "could not type cast taxon")
		}
		res = append(res, t)
	}

	return

}

func (Ω *store) ReadSpecies(cxt context.Context) (res Taxa, err error) {

	q := Ω.FirestoreClient.Collection(CollectionTypeTaxa).
		Where("Rank", "==", RankSpecies).
		Where("Rank", "==", RankSubSpecies).
		Where("Rank", "==", RankForm)

	snapshots, err := q.Documents(cxt).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get taxa")
	}

	for _, doc := range snapshots {
		t := Taxon{}
		if err := doc.DataTo(&t); err != nil {
			return nil, errors.Wrap(err, "could not type cast taxon")
		}
		res = append(res, t)
	}

	return
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
//	if err := Ω.FirebaseClient.Get(context.Background(), txn.Key, &t); err != nil && err != datastore.ErrNoSuchEntity {
//		return nil, errors.Wrap(err, "could not get taxon").SetState(M{utils.LogkeyDatastoreKey: txn.Key.Kind})
//	} else if err == nil {
//		// We have the entity, so exit early.
//		return t.Key, nil
//	}
//
//	k, err := Ω.FirebaseClient.Put(context.Background(), txn.Key, txn)
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
//	_, err := Ω.FirebaseClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
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
//	_, err := Ω.FirebaseClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
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
//	_, err := Ω.FirebaseClient.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
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


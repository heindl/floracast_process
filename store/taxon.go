package store

import (
	"time"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"github.com/saleswise/errors/errors"
	"strings"
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
	RankVariety = TaxonRank("Variety")
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
	"variety": RankVariety,
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
	RankLevelVariety = RankLevel(5)
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
	EcoRegions map[string]int `firestore:",omitempty"`
}

func (Ω Taxon) Combine(s *Taxon) *Taxon {

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

func (Ω *store) CreateTaxonIfNotExists(cxt context.Context, txn Taxon) error {

	// Validate
	if !txn.ID.Valid() {
		return errors.New("invalid taxa")
	}

	//if txn.EcoRegions == nil {
	//	txn.EcoRegions = map[string]int{}
	//}

	ref := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(txn.ID))

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			if strings.Contains(err.Error(), "not found") {
				return tx.Set(ref, txn)
			} else {
				return err
			}
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "could not update occurrence")
	}

	return nil
}

func (Ω *store) SetTaxonPhoto(cxt context.Context, taxonID TaxonID, photoURL string) error {
	if _, err := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(taxonID)).Set(cxt, map[string]interface{}{
		"PhotoURL": photoURL,
	}, firestore.Merge(firestore.FieldPath{"PhotoURL"})); err != nil {
		return errors.Wrap(err, "could not update taxon photo")
	}
	return nil
}
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


// IncrementTaxonEcoRegion updates a map of world wildlife fund eco-regions with the count of occurrences in each.
// This is to be used to sort the taxa to include in each model for the eco region.
// Formula: (OccurrenceCountForTaxonWithinEcoRegion / TotalOccurrenceCountForTaxon) / TotalEcoRegions
// This should prioritize occurrences that only occur within that ecoregion.
func (Ω *store) IncrementTaxonEcoRegion(cxt context.Context, taxonID TaxonID, ecoRegionKey string) error {

	ref := Ω.FirestoreClient.Collection(CollectionTypeTaxa).Doc(string(taxonID))

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(ref)
		if err != nil {
			return errors.Wrapf(err, "could not get taxon: %s", string(taxonID))
		}

		fieldPath := firestore.FieldPath{"EcoRegions", "_"+ecoRegionKey}

		ecoRegionCount, err := doc.DataAtPath(fieldPath)
		if err != nil && !strings.Contains(err.Error(), "no field") && !strings.Contains(err.Error(), `value for field "EcoRegions" is not a map`) {
			return errors.Wrap(err, "could not find region count at field path")
		}

		newCount := int64(1)
		if ecoRegionCount != nil {
			newCount = ecoRegionCount.(int64) + 1
		}
		return tx.Update(ref, []firestore.Update{firestore.Update{FieldPath: fieldPath, Value: newCount}})
	}); err != nil {
		return errors.Wrap(err, "could not update occurrence")
	}

	return nil
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
		Where("RankLevel", "<=", 10)
		//Where("Rank", "==", RankSubSpecies).
		//Where("Rank", "==", RankForm)

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

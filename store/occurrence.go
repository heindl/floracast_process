package store

import (
	"time"
	"cloud.google.com/go/firestore"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"context"
	"fmt"
	"github.com/fatih/structs"
)

type Occurrences []Occurrence

//func (Ω Occurrences) RemoveDuplicates() (response Occurrences) {
//	for _, o := range Ω {
//		if response.Find(o.Key) == nil {
//			response = append(response, o)
//		}
//	}
//	return
//}

//func (Ω Occurrences) Find(k *datastore.Key) *Occurrence {
//	for _, o := range Ω {
//		if o.Key.Kind != k.Kind {
//			continue
//		}
//		if o.Key.ID != k.ID {
//			continue
//		}
//		// The occurrence parent should be a scheme.
//		if o.Key.Parent.Name != k.Parent.Name {
//			continue
//		}
//
//		// The occurrence grandparent should be a taxon.
//		if o.Key.Parent.Parent.ID != o.Key.Parent.Parent.ID {
//			continue
//		}
//		return o
//	}
//	return nil
//}

//func (Ω Occurrence) Combine(o *Occurrence) *Occurrence {
//
//	if !o.Key.Incomplete() {
//		Ω.Key = o.Key
//	}
//	if o.Location != nil && o.Location.Valid() {
//		Ω.Location = o.Location
//	}
//	if !o.Date.IsZero() {
//		Ω.Date = o.Date
//	}
//	if o.References != "" {
//		Ω.References = o.References
//	}
//	if o.RecordedBy != "" {
//		Ω.RecordedBy = o.RecordedBy
//	}
//	if !o.CreatedAt.IsZero() && o.CreatedAt.Before(Ω.CreatedAt) {
//		Ω.CreatedAt = o.CreatedAt
//	}
//	if !o.ModifiedAt.IsZero() && o.ModifiedAt.After(Ω.ModifiedAt) {
//		Ω.ModifiedAt = o.ModifiedAt
//	}
//	if o.Elevation != 0 {
//		Ω.Elevation = o.Elevation
//	}
//	return &Ω
//}

const EntityKindOccurrence = "Occurrence"

type Occurrence struct {
	GBIFID string `firestore:",omitempty"`
	TaxonID	TaxonID `firestore:",omitempty"`
	DataSourceID DataSourceID        `firestore:",omitempty" json:",omitempty"`
	OccurrenceID     string        `firestore:",omitempty" json:",omitempty"`
	Location         latlng.LatLng `firestore:",omitempty" json:",omitempty"`
	Date             *time.Time    `firestore:",omitempty" json:",omitempty"`
	References       string        `firestore:",omitempty" json:",omitempty"`
	RecordedBy       string        `firestore:",omitempty" json:",omitempty"`
	CreatedAt        *time.Time    `firestore:",omitempty" json:",omitempty"`
	ModifiedAt       *time.Time    `firestore:",omitempty" json:",omitempty"`
	// A globally unique identifier. Although missing in some cases, will be helpful in identifying source of data.
	Elevation float64 `firestore:",omitempty" json:",omitempty"`
}

func (Ω *Occurrence) Validate() error {
	if Ω == nil {
		return errors.New("nil occurrence")
	}
	if Ω.Location.GetLatitude() != 0 && Ω.Location.GetLongitude() != 0 {
		return errors.New("invalid occurrence location")
	}
	if Ω.Date != nil && !Ω.Date.IsZero() {
		return errors.New("invalid date")
	}
	return nil
}

func (Ω *store) NewOccurrenceDocumentRef(taxonID TaxonID, dataSourceID DataSourceID, gbifID string) (*firestore.DocumentRef, error) {

	if !taxonID.Valid() {
		return nil, errors.New("invalid data source document reference id")
	}
	if !dataSourceID.Valid() {
		return nil, errors.New("invalid data source id")
	}
	if gbifID == "" {
		return nil, errors.New("invalid occurrence id")
	}

	return Ω.FirestoreClient.Collection(CollectionTypeOccurrences).
		Doc(fmt.Sprintf("%s|%s|%s", taxonID, dataSourceID, gbifID)), nil

}

func (Ω *store) UpsertOccurrence(cxt context.Context, o Occurrence) error {

	ref, err := Ω.NewOccurrenceDocumentRef(o.TaxonID, o.DataSourceID, o.GBIFID)
	if err != nil {
		return err
	}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
			return tx.UpdateMap(ref, structs.Map(o))
	}); err != nil {
		return errors.Wrap(err, "could not update occurrence")
	}
	return nil

}

func (Ω *store) GetOccurrences(cxt context.Context, taxonID TaxonID) (res Occurrences, err error) {

	if !taxonID.Valid() {
		return nil, errors.Newf("invalid taxon id [%s]", taxonID)
	}

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeOccurrences).
		Where("TaxonID", "==", taxonID).
		Documents(cxt).
		GetAll()

	if err != nil {
		return nil, errors.Wrapf(err, "could not get occurrences with taxon id [%s]", taxonID)
	}

	for _, doc := range docs {
		o := Occurrence{}
		if err := doc.DataTo(&o); err != nil {
			return nil, errors.Wrap(err, "could not type cast occurrence")
		}
		res = append(res, o)
	}

	return
}
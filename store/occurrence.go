package store

import (
	"time"
	"cloud.google.com/go/firestore"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"context"
	"fmt"
	"strings"
)

type Occurrences []Occurrence

type Occurrence struct {
	TargetID      string        `firestore:",omitempty" json:",omitempty"`
	TaxonID       TaxonID       `firestore:",omitempty" json:",omitempty"`
	DataSourceID  DataSourceID  `firestore:",omitempty" json:",omitempty"`
	OccurrenceID  string        `firestore:",omitempty" json:",omitempty"`
	Location      latlng.LatLng `firestore:",omitempty" json:",omitempty"`
	Date          *time.Time    `firestore:",omitempty" json:",omitempty"`
	FormattedDate string        `firestore:",omitempty" json:",omitempty"`
	Month         time.Month    `firestore:",omitempty" json:",omitempty"`
	References       string        `firestore:",omitempty" json:",omitempty"`
	RecordedBy       string        `firestore:",omitempty" json:",omitempty"`
	CreatedAt        *time.Time    `firestore:",omitempty" json:",omitempty"`
	ModifiedAt       *time.Time    `firestore:",omitempty" json:",omitempty"`
	// A globally unique identifier. Although missing in some cases, will be helpful in identifying source of data.
	Elevation float64 `firestore:",omitempty" json:",omitempty"`
	EcoRegion string `firestore:",omitempty" json:",omitempty"`
	//S2CellIDs map[string]bool `firestore:",omitempty" json:",omitempty"`
}

var OccurrenceFieldsToMerge = []firestore.FieldPath{
	firestore.FieldPath{"TargetID"},
	firestore.FieldPath{"TaxonID"},
	firestore.FieldPath{"DataSourceID"},
	firestore.FieldPath{"OccurrenceID"},
	firestore.FieldPath{"Location"},
	firestore.FieldPath{"Date"},
	firestore.FieldPath{"FormattedDate"},
	firestore.FieldPath{"Month"},
	firestore.FieldPath{"References"},
	firestore.FieldPath{"RecordedBy"},
	firestore.FieldPath{"ModifiedAt"},
	firestore.FieldPath{"Elevation"},
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

func (Ω *store) NewOccurrenceDocumentRef(taxonID TaxonID, dataSourceID DataSourceID, targetID string) (*firestore.DocumentRef, error) {

	if !taxonID.Valid() {
		return nil, errors.New("invalid data source document reference id")
	}
	if !dataSourceID.Valid() {
		return nil, errors.New("invalid data source id")
	}
	if targetID == "" {
		return nil, errors.New("invalid occurrence id")
	}

	id := fmt.Sprintf("%s|%s|%s", string(taxonID), dataSourceID, targetID)

	return Ω.FirestoreClient.Collection(CollectionTypeOccurrences).Doc(id), nil

}

func (Ω *store) UpsertOccurrence(cxt context.Context, o Occurrence) (isNewOccurrence bool, err error) {

	ref, err := Ω.NewOccurrenceDocumentRef(o.TaxonID, o.DataSourceID, o.TargetID)
	if err != nil {
		return false, err
	}

	isNewOccurrence = false

	//o.S2CellIDs, err = s2Cells(o.Location.Latitude, o.Location.Longitude)
	//if err != nil {
	//	return false, err
	//}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			if strings.Contains(err.Error(), "not found") {
				isNewOccurrence = true
				return tx.Set(ref, o)
			} else {
				return err
			}
		}
		return tx.Set(ref, o, firestore.Merge(OccurrenceFieldsToMerge...))
	}); err != nil {
		return false, errors.Wrap(err, "could not update occurrence")
	}
	return isNewOccurrence, nil
}

//func s2Cells(lat, lng float64) (map[string]bool, error) {
//
//	if lat == 0 || lng == 0 {
//		return nil, errors.New("invalid lat/lng")
//	}
//
//	cell := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
//	cells := map[string]bool{strings.Replace(cell.String(), "/", "_", -1): true}
//	for i:=1; i<14; i++ {
//		cells[strings.Replace(cell.Parent(i).String(), "/", "_", -1)] = true
//	}
//	return cells, nil
//}

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
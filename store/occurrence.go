package store

import (
	"bitbucket.org/heindl/taxa/utils"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"strings"
	"time"
)

type Occurrences []Occurrence

type Occurrence struct {
	TargetID      string             `firestore:",omitempty" json:",omitempty"`
	TaxonID       INaturalistTaxonID `firestore:",omitempty" json:",omitempty"`
	DataSourceID  DataSourceType     `firestore:",omitempty" json:",omitempty"`
	OccurrenceID  string             `firestore:",omitempty" json:",omitempty"`
	Location      latlng.LatLng      `firestore:",omitempty" json:",omitempty"`
	Date          *time.Time         `firestore:",omitempty" json:",omitempty"`
	FormattedDate string             `firestore:",omitempty" json:",omitempty"`
	Month         time.Month         `firestore:",omitempty" json:",omitempty"`
	References    string             `firestore:",omitempty" json:",omitempty"`
	RecordedBy    string        `firestore:",omitempty" json:",omitempty"`
	CreatedAt     *time.Time    `firestore:",omitempty" json:",omitempty"`
	ModifiedAt    *time.Time    `firestore:",omitempty" json:",omitempty"`
	// A globally unique identifier. Although missing in some cases, will be helpful in identifying source of data.
	Elevation   float64 `firestore:",omitempty" json:",omitempty"`
	EcoRegion   string  `firestore:",omitempty" json:",omitempty"`
	CountryCode string
	//S2CellIDs map[string]bool `firestore:",omitempty" json:",omitempty"`
}

func (Ω *Occurrence) mergeFields() []firestore.FieldPath {

	// Required fields
	fields := []firestore.FieldPath{
		{"TargetID"},
		{"INaturalistTaxonID"},
		{"DataSourceID"},
		{"OccurrenceID"},
		{"Location"},
		{"Date"},
		{"FormattedDate"},
		{"Month"},
	}

	if Ω.References != "" {
		fields = append(fields, firestore.FieldPath{"References"})
	}

	if Ω.RecordedBy != "" {
		fields = append(fields, firestore.FieldPath{"RecordedBy"})
	}

	if Ω.Elevation != 0 {
		fields = append(fields, firestore.FieldPath{"Elevation"})
	}

	return fields

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

func (Ω *store) NewOccurrenceDocumentRef(taxonID INaturalistTaxonID, dataSourceID DataSourceType, targetID string) (*firestore.DocumentRef, error) {

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

	isNewOccurrence = false

	ref, err := Ω.NewOccurrenceDocumentRef(o.TaxonID, o.DataSourceID, o.TargetID)
	if err != nil {
		return false, err
	}

	bkf := backoff.NewExponentialBackOff()
	bkf.InitialInterval = time.Second * 1
	ticker := backoff.NewTicker(bkf)
	for _ = range ticker.C {

		//o.S2CellIDs, err = s2Cells(o.Location.Latitude, o.Location.Longitude)
		//if err != nil {
		//	return false, err
		//}

		err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
			if _, err := tx.Get(ref); err != nil {
				if strings.Contains(err.Error(), "not found") {
					isNewOccurrence = true
					return tx.Set(ref, o)
				} else {
					return err
				}
			}
			fields := o.mergeFields()
			return tx.Set(ref, o, firestore.Merge(fields...))
		})

		if err != nil && strings.Contains(err.Error(), "rpc error") {
			fmt.Println("RPC ERROR", err, utils.JsonOrSpew(o))
			continue
		}

		if err != nil {
			ticker.Stop()
			return false, errors.Wrap(err, "could not update occurrence")
		}

		ticker.Stop()
		break
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

func (Ω *store) GetOccurrences(cxt context.Context, taxonID INaturalistTaxonID) (res Occurrences, err error) {

	if !taxonID.Valid() {
		return nil, errors.Newf("invalid taxon id [%s]", taxonID)
	}

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeOccurrences).
		Where("INaturalistTaxonID", "==", taxonID).
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

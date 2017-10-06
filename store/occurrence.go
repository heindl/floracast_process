package store

import (
	"time"
	"cloud.google.com/go/firestore"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"context"
	"fmt"
	"github.com/fatih/structs"
	"strings"
)

type Occurrences []Occurrence

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
		Doc(fmt.Sprintf("%s|%s|%s", string(taxonID), dataSourceID, gbifID)), nil

}

func (Ω *store) UpsertOccurrence(cxt context.Context, o Occurrence) error {

	ref, err := Ω.NewOccurrenceDocumentRef(o.TaxonID, o.DataSourceID, o.GBIFID)
	if err != nil {
		return err
	}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			if strings.Contains(err.Error(), "not found") {
				return tx.Set(ref, o)
			} else {
				return err
			}
		}
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
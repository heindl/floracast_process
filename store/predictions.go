package store

import (
	"context"
	"fmt"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"time"
)

type Prediction struct {
	Location latlng.LatLng `datastore:",omitempty" json:",omitempty"`
	// Date formatted "YYYYMMDD"
	Date                  *time.Time         `datastore:",omitempty" json:",omitempty"`
	CreatedAt             *time.Time         `datastore:",omitempty" json:",omitempty"`
	FormattedDate         string             `datastore:",omitempty" json:",omitempty"`
	Month                 time.Month         `datastore:",omitempty" json:",omitempty"`
	PredictionValue       float64            `datastore:",omitempty" json:",omitempty"`
	ScaledPredictionValue float64            `datastore:",omitempty" json:",omitempty"`
	ScarcityValue         float64            `datastore:",omitempty" json:""`
	TaxonID               INaturalistTaxonID `datastore:",omitempty" json:",omitempty"`
	WildernessAreaName    string             `datastore:",omitempty" json:",omitempty"`
	WildernessAreaID      string             `datastore:",omitempty" json:""`
}

func (Ω *store) PredictionDocumentID(p Prediction) (string, error) {
	if !p.TaxonID.Valid() {
		return "", errors.New("invalid taxon id")
	}
	if p.Date == nil || p.Date.IsZero() {
		return "", errors.New("invalid date")
	}
	if p.Location.GetLatitude() == 0 {
		return "", errors.New("invalid latitude")
	}
	if p.Location.GetLongitude() == 0 {
		return "", errors.New("invalid longitude")
	}
	return fmt.Sprintf("%s|%s|%.6f|%.6f", string(p.TaxonID), p.Date.Format("20060102"), p.Location.GetLatitude(), p.Location.GetLongitude()), nil
}

func (Ω *store) SetPrediction(cxt context.Context, p Prediction) error {
	id, err := Ω.PredictionDocumentID(p)
	if err != nil {
		return err
	}
	if _, err := Ω.FirestoreClient.Collection(CollectionTypePredictions).Doc(id).Set(cxt, p); err != nil {
		return err
	}
	return nil
}

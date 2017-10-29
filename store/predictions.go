package store

import (
	"github.com/saleswise/errors/errors"
	"fmt"
	"google.golang.org/genproto/googleapis/type/latlng"
	"context"
	"time"
)

type Prediction struct{
	Location latlng.LatLng `datastore:",omitempty"`
	// Date formatted "YYYYMMDD"
	Date *time.Time `datastore:",omitempty"`
	CreatedAt *time.Time `datastore:",omitempty"`
	FormattedDate string `datastore:",omitempty"`
	Month time.Month `datastore:",omitempty"`
	PredictionValue float64 `datastore:",omitempty"`
	PercentileOverAllTaxonPredictions float64 `datastore:",omitempty"`
	PercentileOverAllTaxaPredictionsForDay float64 `datastore:",omitempty"`
	TaxonID TaxonID `datastore:",omitempty"`
	WildernessAreaName string `datastore:",omitempty"`
	WildernessAreaID string `datastore:",omitempty"`
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
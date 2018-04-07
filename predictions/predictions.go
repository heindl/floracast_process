package predictions

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/terra/geohashindex"
	"context"
	"github.com/dropbox/godropbox/errors"
)

// Prediction is the standard interface for prediction data.
type Prediction interface {
	NameUsageID() nameusage.ID
	Date() string
	Latitude() float64
	Longitude() float64
	Value() float64
}

// NewPrediction validates and instantiates a new prediction.
func NewPrediction(usageID nameusage.ID, date string, lat, lng, predictionValue float64) (Prediction, error) {
	if !usageID.Valid() {
		return nil, errors.Newf("Prediction requires valid NameUsageID [%s]", usageID)
	}

	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}

	if len(date) != 8 {
		return nil, errors.New("Prediction requires valid Date")
	}

	if predictionValue == 0 {
		return nil, errors.New("Prediction value should be more than 0")
	}

	return &prediction{
		nameUsageID: usageID,
		value:       predictionValue,
		//ScaledPredictionValue: (predictionValue - 0.5) / 0.5,
		formattedDate: date,
		lat:           lat,
		lng:           lng,
	}, nil
}

type prediction struct {
	lat, lng      float64
	nameUsageID   nameusage.ID
	formattedDate string
	value         float64
}

func (Ω *prediction) NameUsageID() nameusage.ID {
	return Ω.nameUsageID
}

func (Ω *prediction) Date() string {
	return Ω.formattedDate
}

func (Ω *prediction) Latitude() float64 {
	return Ω.lat
}

func (Ω *prediction) Longitude() float64 {
	return Ω.lng
}

func (Ω *prediction) Value() float64 {
	return Ω.value
}

// Predictions is a slice of predictions with utility methods.
type Predictions []Prediction

// PredictionWriter is an interface for writing predictions
type PredictionWriter interface {
	WritePredictions(Predictions) error
}

// Upload validates an array of predictions and saves them to FireStore.
func (Ω Predictions) Upload(cxt context.Context, floraStore store.FloraStore) error {

	if len(Ω) == 0 {
		return nil
	}

	c, err := geohashindex.NewCollection(Ω[0].NameUsageID())
	if err != nil {
		return err
	}

	colRef, err := floraStore.FirestoreCollection(store.CollectionGeoIndex)
	if err != nil {
		return err
	}

	for _, p := range Ω {
		if err := c.AddPoint(p.Latitude(), p.Longitude(), p.Date(), p.Value()); err != nil {
			return err
		}
	}

	return c.Upload(cxt, colRef)

}

//func (Ω *store) PredictionDocumentID(p Prediction) (string, error) {
//	if !p.TaxonID.Valid() {
//		return "", errors.New("invalid taxon id")
//	}
//	if p.Date == nil || p.Date.IsZero() {
//		return "", errors.New("invalid date")
//	}
//	if p.Location.GetLatitude() == 0 {
//		return "", errors.New("invalid latitude")
//	}
//	if p.Location.GetLongitude() == 0 {
//		return "", errors.New("invalid longitude")
//	}
//	return fmt.Sprintf("%s|%s|%.6f|%.6f", string(p.TaxonID), p.Date.Format("20060102"), p.Location.GetLatitude(), p.Location.GetLongitude()), nil
//}

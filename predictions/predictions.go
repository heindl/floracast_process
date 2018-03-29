package predictions

import (
	"context"
	"time"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/terra/geoembed"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// Prediction is the standard interface for prediction data.
type Prediction interface {
	UsageID() (nameusage.ID, error)
	Date() (string, error)
	ProtectedArea() (geoembed.CoordinateKey, error)
	ScaledPrediction() (float64, error)
	LatLng() (float64, float64, error)
}

// NewPrediction validates and instantiates a new prediction.
func NewPrediction(usageID nameusage.ID, date string, lat, lng, predictionValue float64) (Prediction, error) {
	if !usageID.Valid() {
		return nil, errors.Newf("Could not create Prediction with invalid ID [%s]", usageID)
	}

	coordinateKey, err := geoembed.NewCoordinateKey(lat, lng)
	if err != nil {
		return nil, err
	}

	if len(date) != 8 {
		return nil, errors.New("Could not create Prediction with Invalid Date")
	}

	if predictionValue == 0 {
		return nil, errors.New("Prediction value should be more than 0")
	}

	return &prediction{
		ProtectedAreaID:       coordinateKey,
		PredictionValue:       predictionValue,
		ScaledPredictionValue: (predictionValue - 0.5) / 0.5,
		FormattedDate:         date,
		GeoPoint: &latlng.LatLng{
			Latitude:  lat,
			Longitude: lng,
		},
	}, nil
}

type prediction struct {
	// Date formatted "YYYYMMDD"
	GeoPoint              *latlng.LatLng `firestore:",omitempty" json:",omitempty"`
	NameUsageID           nameusage.ID   `firestore:",omitempty" json:",omitempty"`
	FormattedDate         string         `firestore:",omitempty" json:",omitempty"`
	Month                 time.Month     `firestore:",omitempty" json:",omitempty"`
	PredictionValue       float64        `firestore:",omitempty" json:",omitempty"`
	ScaledPredictionValue float64        `firestore:",omitempty" json:",omitempty"`
	//ScarcityValue         float64            `firestore:"" json:""`
	//TaxonID               INaturalistTaxonID `datastore:",omitempty" json:",omitempty"`
	ProtectedAreaName string                 `firestore:",omitempty" json:",omitempty"`
	ProtectedAreaSize float64                `firestore:",omitempty" json:",omitempty"`
	ProtectedAreaID   geoembed.CoordinateKey `firestore:"" json:""`
}

func (Ω *prediction) UsageID() (nameusage.ID, error) {
	return Ω.NameUsageID, nil
}
func (Ω *prediction) Date() (string, error) {
	if len(Ω.FormattedDate) != 8 {
		return "", errors.Newf("Invalid Prediction Date [%s]", Ω.FormattedDate)
	}
	return Ω.FormattedDate, nil
}
func (Ω *prediction) ProtectedArea() (geoembed.CoordinateKey, error) {
	return Ω.ProtectedAreaID, nil
}
func (Ω *prediction) ScaledPrediction() (float64, error) {
	return Ω.ScaledPredictionValue, nil
}
func (Ω *prediction) LatLng() (float64, float64, error) {
	return Ω.GeoPoint.GetLatitude(), Ω.GeoPoint.GetLongitude(), nil
}

// Predictions is a slice of predictions with utility methods.
type Predictions []Prediction

// PredictionWriter is an interface for writing predictions
type PredictionWriter interface {
	WritePredictions(Predictions) error
}

// Upload validates an array of predictions and saves them to FireStore.
func (Ω Predictions) Upload(cxt context.Context, writer PredictionWriter) error {

	if writer == nil {
		return errors.New("Valid PredictionWriter required at this point")
	}

	// TODO: Get all ProtectedArea Names.
	// Skip for now ...
	//areaCache, err := protectedarea.NewProtectedAreaCache(florastore)
	//if err != nil {
	//	return err
	//}

	return writer.WritePredictions(Ω)
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

package main

import (
	"time"
	"google.golang.org/appengine"
)

type Prediction struct {
	// Date formatted "YYYYMMDD"
	Location *appengine.GeoPoint `json:",omitempty"`
	Date                  *time.Time         `json:",omitempty"`
	CreatedAt             *time.Time         `json:",omitempty"`
	FormattedDate         string             `json:",omitempty"`
	Month                 time.Month         `json:",omitempty"`
	PredictionValue       float64            `json:",omitempty"`
	ScaledPredictionValue float64            `json:",omitempty"`
	ScarcityValue         float64            `json:""`
	//TaxonID               INaturalistTaxonID `datastore:",omitempty" json:",omitempty"`
	WildernessAreaName    string             `json:",omitempty"`
	WildernessAreaID      string             `json:""`
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
//
//func (Ω *store) SetPrediction(cxt context.Context, p Prediction) error {
//	id, err := Ω.PredictionDocumentID(p)
//	if err != nil {
//		return err
//	}
//	if _, err := Ω.FirestoreClient.Collection(CollectionTypePredictions).Doc(id).Set(cxt, p); err != nil {
//		return err
//	}
//	return nil
//}

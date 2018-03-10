package cache

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
)

// PredictionCache is an interface for reading and writing predictions.
type PredictionCache interface {
	WritePredictions(predictions.Predictions) error
	ReadPredictions(lat, lng, radius float64, qDate string, usageID *nameusage.ID) ([]string, error)
	Close() error
}

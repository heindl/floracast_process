package cache

import (
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/nameusage/nameusage"
)

type PredictionCache interface{
	WritePredictions(predictions.Predictions) error
	ReadPredictions(lat, lng, radius float64, qDate string, usageID *nameusage.NameUsageID) ([]string, error)
	Close() error
}

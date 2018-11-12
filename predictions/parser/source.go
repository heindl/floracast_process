package parser

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"context"
)

const GCSPredictionsPath = "predictions"

// PredictionResult is a standard struct for prediction data.
type PredictionResult struct {
	Latitude, Longitude float64
	Date                string
	Target, Random      float64
	NameUsageID         nameusage.ID
}

// PredictionSource is an interface provider to sources of prediction data.
type PredictionSource interface {
	FetchLatestPredictionFileNames(cxt context.Context, id nameusage.ID, date string) ([]string, error)
	FetchPredictions(cxt context.Context, gcsFilePath string) ([]*PredictionResult, error)
}

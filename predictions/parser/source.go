package parser

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"context"
)

const GCSPredictionsPath = "predictions"

type PredictionResult struct {
	Latitude, Longitude float64
	Date                string
	Target, Random      float64
	NameUsageID         nameusage.ID
}

type PredictionSource interface {
	FetchLatestPredictionFileNames(cxt context.Context, id nameusage.ID, date string) ([]string, error)
	FetchPredictions(cxt context.Context, gcsFilePath string) ([]*PredictionResult, error)
}

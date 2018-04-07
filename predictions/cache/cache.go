package cache

import (
	"bitbucket.org/heindl/process/predictions"
)

// PredictionCache is an interface for reading and writing predictions.
type PredictionCache interface {
	WritePredictions(predictions.Predictions) error
	ReadPredictions(bboxString string) (predictions.Predictions, error)
	Close() error
}

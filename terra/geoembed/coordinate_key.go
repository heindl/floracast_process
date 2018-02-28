package geoembed

import (
	"bitbucket.org/heindl/process/terra/geo"
	"fmt"
	"strings"
)

type CoordinateKey string

func (Ω CoordinateKey) Valid() bool {
	return len(Ω) == 11
}

func NewCoordinateKey(lat, lng float64) (CoordinateKey, error) {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return "", err
	}
	// Intentionally reduce the precision of the coordinates to ensure we're not duplicating occurrences.
	k := fmt.Sprintf("%.3f_%.3f", lat, lng)
	// Replace decimals in order to use as firestore key.
	k = strings.Replace(k, ".", "|", -1)
	return CoordinateKey(k), nil
}

func (Ω *GeoFeatureSet) CoordinateKey() (CoordinateKey, error) {
	return NewCoordinateKey(Ω.Lat(), Ω.Lng())
}

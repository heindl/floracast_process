package geoembed

import (
	"bitbucket.org/heindl/process/terra/geo"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

// CoordinateKey is lat/lng formatted string to the third decimal that can be used as a Firestore ID
// The original intention was to ensure no occurrence duplications exist
type CoordinateKey string

// Valid checks a CoordinateKey for expected string length
func (Ω CoordinateKey) Valid() bool {
	return len(Ω) == 14
}

// NewCoordinateKey validates a lat/lng to be in expected bounds, and returns a new string key
func NewCoordinateKey(lat, lng float64) (CoordinateKey, error) {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return "", errors.Wrapf(err, "Invalid Coordinates for CoordinateKey [%f, %f]", lat, lng)
	}
	// Intentionally reduce the precision of the coordinates to ensure we're not duplicating occurrences.
	k := fmt.Sprintf("%.3f_%.3f", lat, lng)
	// Replace decimals in order to use as firestore key.
	k = strings.Replace(k, ".", "|", -1)
	return CoordinateKey(k), nil
}

// CoordinateKey generates a new key from a GeoFeatureSet latlng.LatLng
func (Ω *GeoFeatureSet) CoordinateKey() (CoordinateKey, error) {
	return NewCoordinateKey(Ω.Lat(), Ω.Lng())
}

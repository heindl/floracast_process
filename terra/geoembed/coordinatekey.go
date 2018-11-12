package geoembed

import (
	"github.com/heindl/floracast_process/terra/geo"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
)

// CoordinateKey is lat/lng formatted string to the third decimal that can be used as a Firestore ID
// The original intention was to ensure no occurrence duplications exist
//type CoordinateKey string

type S2Key string

func (Ω S2Key) Valid() bool {
	return len(Ω) == 9
}

func NewS2Key(lat, lng float64) (S2Key, error) {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return "", errors.Wrapf(err, "Invalid Coordinates for CoordinateKey [%f, %f]", lat, lng)
	}
	return S2Key(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)).Parent(16).ToToken()), nil
}

func (Ω S2Key) Parse() (float64, float64, error) {

	if !Ω.Valid() {
		return 0, 0, errors.New("Invalid S2Key")
	}

	cellID := s2.CellIDFromToken(string(Ω))

	lat := cellID.LatLng().Lat.Degrees()
	lng := cellID.LatLng().Lng.Degrees()

	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return 0, 0, errors.Wrapf(err, "Invalid Coordinates for CoordinateKey [%f, %f]", lat, lng)
	}

	return lat, lng, nil
}

// Valid checks a CoordinateKey for expected string length
//func (Ω CoordinateKey) Valid() bool {
//	return len(Ω) == 14
//}
//
//func (Ω CoordinateKey) Parse() (*latlng.LatLng, error) {
//	if !Ω.Valid() {
//		return nil, errors.Newf("Invalid CoordinateKey [%s]", string(Ω))
//	}
//	s := strings.Split(strings.Replace(string(Ω), "|", ".", -1), "_")
//	lat, err := strconv.ParseFloat(s[0], 10)
//	if err != nil {
//		return nil, err
//	}
//	lng, err := strconv.ParseFloat(s[1], 10)
//	if err != nil {
//		return nil, err
//	}
//	return &latlng.LatLng{lat, lng}, nil
//}
//
//// NewCoordinateKey validates a lat/lng to be in expected bounds, and returns a new string key
//func NewCoordinateKey(lat, lng float64) (CoordinateKey, error) {
//	if err := geo.ValidateCoordinates(lat, lng); err != nil {
//		return "", errors.Wrapf(err, "Invalid Coordinates for CoordinateKey [%f, %f]", lat, lng)
//	}
//	// Intentionally reduce the precision of the coordinates to ensure we're not duplicating occurrences.
//	k := fmt.Sprintf("%.3f_%.3f", lat, lng)
//	// Replace decimals in order to use as firestore key.
//	k = strings.Replace(k, ".", "|", -1)
//	return CoordinateKey(k), nil
//}

// CoordinateKey generates a new key from a GeoFeatureSet latlng.Coordinates
//func (Ω *GeoFeatureSet) CoordinateKey() (CoordinateKey, error) {
//	return NewCoordinateKey(Ω.Lat(), Ω.Lng())
//}

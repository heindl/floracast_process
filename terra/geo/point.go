package geo

import (
	"errors"
	dropboxErrors "github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
	pmgeo "github.com/paulmach/go.geo"
	"math"
)

type Point struct {
	latlng     *s2.LatLng
	properties map[string]interface{}
}

var EPSILON float64 = 0.00001

func CoordinatesEqual(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

// Error response when coordinates are invalid.
var ErrInvalidCoordinates = errors.New("Invalid Coordinates")

// Utility function for ensuring coordinates fall within North America and has sufficient resolution.
func ValidateCoordinates(lat, lng float64) error {

	if lat < 6.6 || lat > 83.3 {
		return dropboxErrors.Wrapf(ErrInvalidCoordinates, "Latitude [%f] is out of bounds", lat)
	}
	if lng < -178.2 || lng > -49.0 {
		return dropboxErrors.Wrapf(ErrInvalidCoordinates, "Longitude [%f] is out of bounds", lng)
	}
	// We need the decimal precision to be at least a football field, so require at least three decimal places (110m).
	if hasDecimalPlaces(2, lat) || hasDecimalPlaces(1, lat) {
		return dropboxErrors.Wrapf(ErrInvalidCoordinates, "Latitude [%f] has insufficient precision", lat)
	}
	if hasDecimalPlaces(2, lng) || hasDecimalPlaces(1, lng) {
		return dropboxErrors.Wrapf(ErrInvalidCoordinates, "Longitude [%f] has insufficient precision", lng)
	}
	return nil
}

func hasDecimalPlaces(i int, v float64) bool {
	vf := v * math.Pow(10.0, float64(i))
	extra := vf - float64(int(vf))
	return extra == 0
}

func NewPoint(lat, lng float64) (*Point, error) {
	ll := s2.LatLngFromDegrees(lat, lng)
	if !ll.IsValid() {
		return nil, errors.New("could not create new point with invalid latlng")
	}

	p := Point{
		latlng:     &ll,
		properties: map[string]interface{}{},
	}

	if p.IsZero() {
		return nil, errors.New("new point is zero")
	}

	return &p, nil
}

func (Ω *Point) IsZero() bool {
	if Ω.latlng == nil || Ω.Latitude() == 0 || Ω.Longitude() == 0 {
		return true
	}
	return false
}

func (Ω Point) Latitude() float64 {
	return Ω.latlng.Lat.Degrees()
}

func (Ω *Point) S2TokenMap() map[int]string {
	initial_cell_id := s2.CellIDFromLatLng(*Ω.latlng)
	feature_array := map[int]string{}
	for i := 0; i < 8; i++ {
		cell_id := initial_cell_id.Parent(i)
		feature_array[i] = cell_id.ToToken()
	}
	return feature_array
}

func (Ω Point) Longitude() float64 {
	return Ω.latlng.Lng.Degrees()
}

func (Ω *Point) SetProperty(key string, value interface{}) error {
	if key == "" {
		return errors.New("Invalid Point property key")
	}
	if value == nil {
		return errors.New("Invalid Point property value")
	}
	Ω.properties[key] = value
	return nil
}

func (Ω *Point) AsArray() []float64 {
	return []float64{Ω.Longitude(), Ω.Latitude()}
}

func (Ω *Point) DistanceKilometers(np *Point) float64 {
	p1 := pmgeo.NewPointFromLatLng(Ω.Latitude(), Ω.Longitude())
	p2 := pmgeo.NewPointFromLatLng(np.Latitude(), np.Longitude())
	return p1.GeoDistanceFrom(p2) / 1000.0
}

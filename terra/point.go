package terra

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
	pmgeo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geojson"
)

type Point struct {
	latlng *s2.LatLng
	properties map[string]interface{}
}

var EPSILON float64 = 0.00001

func CoordinatesEqual(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func NewPoint(lat, lng float64) (*Point, error) {
	ll := s2.LatLngFromDegrees(lat, lng)
	if !ll.IsValid() {
		return nil, errors.New("could not create new point with invalid latlng")
	}

	p := Point{
		latlng: &ll,
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

type Points []*Point

func (Ω Points) Centroid() (*Point, error) {
	pointset := &pmgeo.PointSet{}
	for _, p := range Ω {
		pointset = pointset.Push(pmgeo.NewPointFromLatLng(p.Latitude(), p.Longitude()))
	}
	centroid := pointset.GeoCentroid()
	return NewPoint(centroid.Lat(), centroid.Lng())
}

func (Ω Points) GeoJSON() ([]byte, error) {
	fc := geojson.NewFeatureCollection()
	for _, pt := range Ω {
		f := geojson.NewPointFeature(pt.AsArray())
		for k, v := range pt.properties {
			f.SetProperty(k, v)
		}
		fc = fc.AddFeature(f)
	}
	b, err := fc.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal points")
	}
	return b, nil
}

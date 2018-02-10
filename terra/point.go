package terra

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
	pmgeo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geojson"
)

type Point struct {
	latlng *s2.LatLng
}



var EPSILON float64 = 0.00001

func CoordinatesEqual(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func NewPoint(lat, lng float64) Point {
	ll := s2.LatLngFromDegrees(lat, lng)
	return Point{&ll}
}

func (Ω Point) IsZero() bool {
	if Ω.latlng == nil || Ω.Latitude() == 0 || Ω.Longitude() == 0 {
		return true
	}
	return false
}

func (Ω Point) Latitude() float64 {
	return Ω.latlng.Lat.Degrees()
}

func (Ω Point) S2TokenMap() map[string]bool {
	initial_cell_id := s2.CellIDFromLatLng(*Ω.latlng)
	feature_array := map[string]bool{}
	for i := 0; i < 10; i++ {
		cell_id := initial_cell_id.Parent(i)
		feature_array[cell_id.ToToken()] = true
	}
	return feature_array
}

func (Ω Point) Longitude() float64 {
	return Ω.latlng.Lng.Degrees()
}

func (Ω Point) AsArray() []float64 {
	return []float64{Ω.Longitude(), Ω.Latitude()}
}

func (Ω Point) DistanceKilometers(np Point) float64 {
	p1 := pmgeo.NewPointFromLatLng(Ω.Latitude(), Ω.Longitude())
	p2 := pmgeo.NewPointFromLatLng(np.Latitude(), np.Longitude())
	return p1.GeoDistanceFrom(p2) / 1000.0
}

type Points []Point

func (Ω Points) Centroid() Point {
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
		fc = fc.AddFeature(geojson.NewPointFeature(pt.AsArray()))
	}
	b, err := fc.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal points")
	}
	return b, nil
}

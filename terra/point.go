package terra

import (
	"github.com/golang/geo/s2"
	"github.com/paulmach/go.geojson"
	pmgeo "github.com/paulmach/go.geo"
	"github.com/dropbox/godropbox/errors"
)

type Point struct {
	latlng *s2.LatLng
}

func NewPoint(lat, lng float64) Point {
	ll := s2.LatLngFromDegrees(lat, lng)
	return Point{&ll}
}

func (Ω Point) Empty() bool {
	if Ω.latlng == nil || Ω.Latitude() == 0 {
		return true
	}
	return false
}

func (Ω Point) Latitude() float64 {
	return Ω.latlng.Lat.Degrees()
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
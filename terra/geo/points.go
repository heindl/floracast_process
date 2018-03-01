package geo

import (
	"github.com/dropbox/godropbox/errors"
	pmgeo "github.com/paulmach/go.geo"
	"github.com/paulmach/go.geojson"
)

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

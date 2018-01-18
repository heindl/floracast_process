package terra

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/paulmach/go.geojson"
	"github.com/mongodb/mongo-tools/common/json"
	"github.com/buger/jsonparser"
)

type Feature struct {
	polyLabel Point
	area float64
	properties []byte
	multiPolygon MultiPolygon
}

func (Ω *Feature) Normalize() {
	Ω.polyLabel = Ω.multiPolygon.PolylabelOfLargestPolygon()
	Ω.area = Ω.multiPolygon.Area()
}

func (Ω *Feature) SetProperties(b []byte) {
	Ω.properties = b
}

func (Ω *Feature) GetProperties(i interface{}) error {
	if err := json.Unmarshal(Ω.properties, i); err != nil {
		return errors.Wrap(err, "could not unmarshal properties")
	}
	return nil
}

func (Ω *Feature) PushMultiPolygon(m MultiPolygon) {
	Ω.multiPolygon = Ω.multiPolygon.PushMultiPolygon(m)
	Ω.Normalize()
}

func (Ω *Feature) Valid() bool {
	if Ω.multiPolygon.Empty() {
		return false
	}
	return true
}

func (Ω *Feature) PolyLabel() Point {
	return Ω.polyLabel
}

func (Ω *Feature) Area() float64 {
	return Ω.area
}

func (Ω *Feature) Contains(lat, lng float64) bool {
	return Ω.multiPolygon.Contains(lat, lng)
}

func (Ω *Feature) GeoJSON() ([]byte, error) {

	feature := geojson.NewFeature(geojson.NewMultiPolygonGeometry(Ω.multiPolygon.ToArray()...))

	if Ω.properties != nil {
		m := map[string]interface{}{}
		if err := json.Unmarshal(Ω.properties, &m); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal properties to map")
		}
		for k, v := range m {
			feature.SetProperty(k, v)
		}
	}
	feature.SetProperty("terra_area", Ω.Area())
	feature.SetProperty("terra_polylabel", Ω.polyLabel.AsArray())

	encoded_feature, err := feature.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal geojson feature")
	}

	return encoded_feature, nil
}

func (Ω *Feature) GetProperty(prop string) (interface{}, error) {
	v, _, _, err := jsonparser.Get(Ω.properties, prop)
	if err != nil {
		return nil, errors.Wrap(err, "could not get value")
	}
	return v, nil
}

func (Ω *Feature) GetPropertyString(prop string) (string, error) {
	s, err := jsonparser.GetString(Ω.properties, prop)
	if err != nil {
		return "", errors.Wrap(err, "could not get string")
	}
	return s, nil
}
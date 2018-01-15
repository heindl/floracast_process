package terra

import (
	"io/ioutil"
	"bufio"
	"os"
	"github.com/buger/jsonparser"
	"gopkg.in/tomb.v2"
	"github.com/dropbox/godropbox/errors"
	"github.com/paulmach/go.geojson"
	"github.com/golang/geo/s2"
	"encoding/json"
)

type GeoJSONParsedCallback func(encoded_properties []byte, polygon MultiPolygon) error

func ReadGeoJSONFeatureCollectionFile(filepath string, callback GeoJSONParsedCallback) error {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}
	return ParseGeoJSONFeatureCollection(b, callback)
}

func ParseGeoJSONFeatureCollection(encodedFeatureCollection []byte, callback GeoJSONParsedCallback) error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		if _, err := jsonparser.ArrayEach(encodedFeatureCollection, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			tmb.Go(func() error {
				if err != nil {
					return err
				}
				return ParseGeoJSONFeature(value, callback)
			})
		}, "features"); err != nil {
			return errors.Wrap(err, "could not parse features array")
		}
		return nil
	})
	return tmb.Wait();
}

func ParseGeoJSONFeature(encoded_feature []byte, callback GeoJSONParsedCallback) error {

	endoded_properties, _, _, err := jsonparser.Get(encoded_feature, "properties")
	if err != nil {
		return errors.Wrap(err, "could not get properties")
	}

	encoded_geometry, _, _, err := jsonparser.Get(encoded_feature, "geometry")
	if err != nil {
		return errors.Wrap(err, "could not get geometry")
	}

	geometry, err := geojson.UnmarshalGeometry(encoded_geometry)
	if err != nil {
		return errors.Wrap(err, "could could not unmarshal geometry")
	}

	if !geometry.IsMultiPolygon() && !geometry.IsPolygon() {
		return errors.Newf("unsupported geometry type: %s", geometry.Type)
	}
	var multipolygon MultiPolygon
	if geometry.IsMultiPolygon() {
		// MultiPolygon    [][][][]float64
		for _, polygon_array := range geometry.MultiPolygon {
			np, err := NewPolygon(polygon_array)
			if err != nil {
				return err
			}
			multipolygon = multipolygon.PushPolygon(np)
		}
	} else {
		// Polygon         [][][]float64
		np, err := NewPolygon(geometry.Polygon)
		if err != nil {
			return err
		}
		multipolygon = multipolygon.PushPolygon(np)
	}

	return callback(endoded_properties, multipolygon)
}

func (Ω MultiPolygon) ToGeoJSONFeatureCollection(encoded_properties []byte) ([]byte, error) {
	encoded_feature, err := Ω.ToGeoJSONFeature(encoded_properties)
	if err != nil {
		return nil, err
	}
	feature, err := geojson.UnmarshalFeature(encoded_feature)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal geojson feature")
	}

	fc := geojson.NewFeatureCollection()
	fc.AddFeature(feature)
	encoded_feature_collection, err := fc.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal geojson feature collection")
	}
	return encoded_feature_collection, nil
}

func (Ω MultiPolygon) ToGeoJSONFeature(encoded_properties []byte) ([]byte, error) {

	multipolygon_array := [][][][]float64{}
	for _, polygon := range Ω {
		polygon_array := [][][]float64{}
		for _, loop := range polygon.Loops() {
			loop_array := [][]float64{}
			// Reverse the order for geojson format.
			vertices := loop.Vertices()
			for i := len(vertices)-1; i >= 0; i-- {
				coords := s2.LatLngFromPoint(vertices[i])
				loop_array = append(loop_array, []float64{coords.Lng.Degrees(), coords.Lat.Degrees()})
			}
			polygon_array = append(polygon_array, loop_array)
		}
		multipolygon_array = append(multipolygon_array, polygon_array)
	}

	feature := geojson.NewFeature(geojson.NewMultiPolygonGeometry(multipolygon_array...))

	if encoded_properties != nil {
		m := map[string]interface{}{}
		if err := json.Unmarshal(encoded_properties, &m); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal properties to map")
		}
		for k, v := range m {
			feature.SetProperty(k, v)
		}
	}

	encoded_feature, err := feature.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal geojson feature")
	}

	return encoded_feature, nil
}
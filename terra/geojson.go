package terra

import (
	"bufio"
	"github.com/buger/jsonparser"
	"github.com/dropbox/godropbox/errors"
	"github.com/paulmach/go.geojson"
	"gopkg.in/tomb.v2"
	"io/ioutil"
	"os"
	"sync"
)

type GeoJSONParsedCallback func(feature *Feature) error

func ReadFeatureCollectionFromGeoJSONFile(filepath string, property_filter func([]byte) bool) (*FeatureCollection, error) {

	fc_holder := []*Feature{}
	m := sync.Mutex{}

	callback := func(f *Feature) error {
		m.Lock()
		defer m.Unlock()
		if property_filter != nil && property_filter(f.properties) {
			return nil
		}
		fc_holder = append(fc_holder, f)
		return nil
	}

	if err := ReadFeaturesFromGeoJSONFeatureCollectionFile(filepath, callback); err != nil {
		return nil, err
	}

	fc := FeatureCollection{}
	if err := fc.Append(fc_holder...); err != nil {
		return nil, err
	}

	return &fc, nil
}

func ReadFeaturesFromGeoJSONFeatureCollectionFile(filepath string, callback GeoJSONParsedCallback) error {
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
					return errors.Wrapf(err, "could not parse features array: %s", string(value))
				}
				if err := ParseGeoJSONFeature(value, callback); err != nil {
					return errors.Wrapf(err, "could not parse features array: %s", string(value))
				}
				return nil
			})
		}, "features"); err != nil {
			return err
		}
		return nil
	})
	return tmb.Wait()
}

func ParseGeoJSONFeature(encoded_feature []byte, callback GeoJSONParsedCallback) error {

	encoded_properties, _, _, err := jsonparser.Get(encoded_feature, "properties")
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
			multipolygon, err = multipolygon.PushPolygon(np)
			if err != nil {
				return err
			}
		}
	} else {
		// Polygon         [][][]float64
		np, err := NewPolygon(geometry.Polygon)
		if err != nil {
			return err
		}
		multipolygon, err = multipolygon.PushPolygon(np)
		if err != nil {
			return err
		}
	}

	f := Feature{}
	f.PushMultiPolygon(multipolygon)
	f.SetProperties(encoded_properties)

	return callback(&f)
}

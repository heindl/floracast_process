package geo

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

func ReadFeaturesFromGeoJSONFeatureCollectionFile(filepath string, callback GeoJSONParsedCallback) (err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return errors.Wrapf(err, "Could not open file [%s]", filepath)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
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

	multipolygon, err := unmarshalGeometryAsMultiPolygon(encoded_geometry)
	if err != nil {
		return err
	}

	f := Feature{}
	if err := f.PushMultiPolygon(multipolygon); err != nil {
		return err
	}
	f.SetProperties(encoded_properties)

	return callback(&f)
}

func unmarshalGeometryAsMultiPolygon(encoded_geometry []byte) (MultiPolygon, error) {
	geometry, err := geojson.UnmarshalGeometry(encoded_geometry)
	if err != nil {
		return nil, errors.Wrap(err, "could could not unmarshal geometry")
	}

	if !geometry.IsMultiPolygon() && !geometry.IsPolygon() {
		return nil, errors.Newf("unsupported geometry type: %s", geometry.Type)
	}
	if geometry.IsMultiPolygon() {
		// MultiPolygon    [][][][]float64
		mp := MultiPolygon{}
		for _, polygon_array := range geometry.MultiPolygon {
			np, err := NewPolygon(polygon_array)
			if err != nil {
				return nil, err
			}
			mp, err = mp.PushPolygon(np)
			if err != nil {
				return nil, err
			}
		}
		return mp, nil
	}
	// Polygon         [][][]float64
	np, err := NewPolygon(geometry.Polygon)
	if err != nil {
		return nil, err
	}
	return MultiPolygon{}.PushPolygon(np)
}

package geo

import (
	"bitbucket.org/heindl/process/utils"
	"bufio"
	"github.com/buger/jsonparser"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
	"github.com/paulmach/go.geojson"
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
	defer utils.SafeClose(f, &err)

	b, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return err
	}
	return ParseGeoJSONFeatureCollection(b, callback)
}

func ParseGeoJSONFeatureCollection(encodedFeatureCollection []byte, callback GeoJSONParsedCallback) error {

	featureBytes := [][]byte{}
	if _, err := jsonparser.ArrayEach(encodedFeatureCollection, func(f []byte, _ jsonparser.ValueType, _ int, _ error) {
		featureBytes = append(featureBytes, f)
	}, "features"); err != nil {
		return errors.Wrap(err, "Unable to Encoded FeatureCollection")
	}

	//tmb := tomb.Tomb{}
	//tmb.Go(func() error {
	for _, fb := range featureBytes {
		//_fb := 𝝨
		//tmb.Go(func() error {
		//	fb := _fb
		if err := ParseGeoJSONFeature(fb, callback); err != nil {
			return err
		}
		//})
	}
	return nil
	//})
	//return tmb.Wait()
}

func ParseGeoJSONFeature(encodedFeature []byte, callback GeoJSONParsedCallback) error {

	encodedProperties, _, _, err := jsonparser.Get(encodedFeature, "properties")
	if err != nil {
		return errors.Wrap(err, "could not get properties")
	}

	encodedGeometry, _, _, err := jsonparser.Get(encodedFeature, "geometry")
	if err != nil {
		return errors.Wrap(err, "could not get geometry")
	}

	multipolygon, err := unmarshalGeometryAsMultiPolygon(encodedGeometry)
	if err != nil {
		return err
	}

	f := Feature{}
	if err := f.PushMultiPolygon(multipolygon); err != nil {
		return err
	}
	f.SetProperties(encodedProperties)

	return callback(&f)
}

func unmarshalGeometryAsMultiPolygon(encodedGeometry []byte) (MultiPolygon, error) {
	geometry, err := geojson.UnmarshalGeometry(encodedGeometry)
	if err != nil {
		return nil, errors.Wrap(err, "could could not unmarshal geometry")
	}

	if !geometry.IsMultiPolygon() && !geometry.IsPolygon() {
		return nil, errors.Newf("unsupported geometry type: %s", geometry.Type)
	}
	if geometry.IsMultiPolygon() {
		// MultiPolygon    [][][][]float64
		mp := MultiPolygon{}
		for _, polygonArray := range geometry.MultiPolygon {
			np := &s2.Polygon{}
			np, err = NewPolygon(polygonArray)
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

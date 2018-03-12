package tfrecords

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/tfutils"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Record interface {
	GetFeatureString(featureName string, fType FeatureType) (string, error)
	Tensor() (*tf.Tensor, error)
	String() string
	Latitude() (float64, error)
	Longitude() (float64, error)
}

func NewRecord(b []byte) (Record, error) {
	return &record{b: b}, nil
}

type record struct {
	b []byte
}

func (Ω *record) Tensor() (*tf.Tensor, error) {
	t, err := tf.NewTensor([]string{string(Ω.b)})
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate Tensor")
	}
	return t, nil
}

func (Ω *record) String() string {
	return string(Ω.b)
}

func (Ω *record) Latitude() (float64, error) {
	list, err := Ω.GetFloatFeature("latitude")
	if err != nil {
		return 0, nil
	}
	if len(list) == 0 {
		return 0, errors.New("Latitude missing")
	}
	return float64(list[0]), nil
}

func (Ω *record) Longitude() (float64, error) {
	list, err := Ω.GetFloatFeature("longitude")
	if err != nil {
		return 0, nil
	}
	if len(list) == 0 {
		return 0, errors.New("Longitude missing")
	}
	return float64(list[0]), nil
}

func (Ω *record) GetFeatureString(featureName string, fType FeatureType) (string, error) {
	features, err := tfutils.GetFeatureMapFromTFRecord(Ω.b)
	if err != nil {
		return "", errors.Wrap(err, "Could not get features")
	}
	m := features.GetFeature()
	f, ok := m[featureName]
	if !ok {
		return "", errors.Newf("Feature [%s] missing from TFRecord", featureName)
	}
	var res fmt.Stringer
	switch fType {
	case FeatureTypeBytes:
		res = f.GetBytesList()
	case FeatureTypeInt:
		res = f.GetInt64List()
	case FeatureTypeFloat:
		res = f.GetFloatList()
	default:
		return "", errors.Newf("Unsupported FeatureType [%s]", fType)
	}
	return res.String(), nil
}

func (Ω *record) GetFloatFeature(featureName string) ([]float32, error) {
	features, err := tfutils.GetFeatureMapFromTFRecord(Ω.b)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get features")
	}
	m := features.GetFeature()
	f, ok := m[featureName]
	if !ok {
		return nil, errors.Newf("Feature [%s] missing from TFRecord", featureName)
	}
	list := f.GetFloatList()
	return list.Value, nil
}

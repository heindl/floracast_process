package main

import (
	"bitbucket.org/heindl/process/store"
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/tfutils"
	"io"
	"strings"
)

func newTFRecorder(ctx context.Context, floraStore store.FloraStore, gcsFilePath string) (*tfRecorder, error) {
	gcsObjects, err := floraStore.CloudStorageObjects(ctx, gcsFilePath, ".tfrecords")
	if err != nil {
		return nil, err
	}
	return &tfRecorder{
		gcsObjects: gcsObjects,
	}, nil
}

type tfRecorder struct {
	gcsObjects []*storage.ObjectHandle
}

type tfEachCallback func([]byte) error

func (Ω *tfRecorder) CountRecords(ctx context.Context) (int, error) {
	count := 0
	err := Ω.Each(ctx, func(_ []byte) error {
		count++
		return nil
	})
	return count, err
}

type featureType int

const (
	featureTypeBytes featureType = iota + 1
	featureTypeInt
	featureTypeFloat
)

func (Ω *tfRecorder) PrintFeature(ctx context.Context, featureName string, fType featureType) error {
	return Ω.Each(ctx, func(b []byte) error {
		features, err := tfutils.GetFeatureMapFromTFRecord(b)
		if err != nil {
			return errors.Wrap(err, "Could not get features")
		}
		m := features.GetFeature()
		f, ok := m[featureName]
		if !ok {
			return errors.Newf("Feature [%s] missing from TFRecord", featureName)
		}
		var res fmt.Stringer
		switch fType {
		case featureTypeBytes:
			res = f.GetBytesList()
		case featureTypeInt:
			res = f.GetInt64List()
		case featureTypeFloat:
			res = f.GetFloatList()
		default:
			return errors.Newf("Unsupported FeatureType [%s]", fType)
		}
		fmt.Println(res.String())
		return nil
	})
}

func (Ω *tfRecorder) Each(ctx context.Context, cb tfEachCallback) error {
	for _, gcsObject := range Ω.gcsObjects {
		if err := Ω.callObject(ctx, gcsObject, cb); err != nil {
			return err
		}
	}
	return nil
}

func (Ω *tfRecorder) callObject(ctx context.Context, gcsObject *storage.ObjectHandle, cb tfEachCallback) error {
	// This slows things down, but to be sure of data consistency ...
	attrs, err := gcsObject.Attrs(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(attrs.Name, ".tfrecords") {
		return errors.Newf("Expected a tfrecords encoded file [%s]", attrs.Name)
	}
	//
	//if attrs.ContentEncoding != "gzip" {
	//	if _, err := gcsObject.Update(ctx, storage.ObjectAttrsToUpdate{ContentEncoding: "gzip"}); err != nil {
	//		return 0, errors.Wrapf(err, "Could not update ContentEncoding of GCS Object [%s]", attrs.Name)
	//	}
	//}

	gcsObjectReader, err := gcsObject.ReadCompressed(false).NewReader(ctx)
	if err != nil {
		return err
	}

	reader, err := tfutils.NewReaderFromBufio(bufio.NewReader(gcsObjectReader), &tfutils.RecordReaderOptions{
		CompressionType: tfutils.CompressionTypeNone,
	})
	if err != nil {
		return errors.Wrap(err, "Could not get TfRecordReader")
	}
	for {
		r, err := reader.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "Could not read TfRecord")
		}
		if err := cb(r); err != nil {
			return err
		}
	}
	return nil
}
